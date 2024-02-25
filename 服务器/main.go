package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recoverMiddleware "github.com/gofiber/fiber/v2/middleware/recover"

	"2024challenge/dynamicflag"
	"2024challenge/game2048"
	"2024challenge/yolov5verify"
)

func main() {
	var listenAddr string
	flag.StringVar(&listenAddr, "listen", ":3000", "HTTP 监听 IP:端口号")
	flag.Parse()

	app := fiber.New()
	app.Use(recoverMiddleware.New(recoverMiddleware.Config{
		EnableStackTrace: true,
	}))
	app.Use(logger.New(logger.Config{
		Format:        "${ip}:${port} - ${reqHeader:X-Real-Ip} - [uid:${cookie:uid}] [${time}] \"${method} ${url}\" ${status} ${latency} ${bytesSent} \"${error}\" \"${referer}\" \"${ua}\"\n",
		TimeFormat:    time.RFC3339Nano,
		TimeZone:      "Asia/Shanghai",
		DisableColors: true,
	}))
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: base64.StdEncoding.EncodeToString([]byte("jSIP3U61cbDsTHwRyn7ayi2ttuI37QA9")),
	}))

	app.Post("/auth/login", func(c *fiber.Ctx) error {
		uidStr := c.FormValue("uid")
		uid, _ := strconv.ParseUint(uidStr, 10, 64)
		if uid == 0 {
			c.Status(400)
			return c.SendString("uid 必须为正整数")
		}
		c.Cookie(&fiber.Cookie{
			Name:  "uid",
			Value: strconv.FormatUint(uid, 10),
		})
		flagAContent, expiredAt := dynamicflag.CalcFlag(strconv.FormatUint(uid, 10), "_2024_52pojie_flag_aaa_", time.Now())
		c.Cookie(&fiber.Cookie{
			Name:    "flagA",
			Value:   fmt.Sprintf("flagA{%s}", flagAContent),
			Expires: expiredAt,
		})
		return c.Redirect("/")
	})
	app.Post("/auth/logout", func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:     "uid",
			Expires:  time.Now().Add(-(time.Hour * 2)),
			HTTPOnly: true,
			SameSite: "lax",
		})
		return c.Redirect("/")
	})
	app.Get("/auth/uid", func(c *fiber.Ctx) error {
		return c.SendString(c.Cookies("uid"))
	})

	app.Post("/flagC/verify", func(c *fiber.Ctx) error {
		uidStr := c.Cookies("uid")
		uid, _ := strconv.ParseUint(uidStr, 10, 64)
		if uid == 0 {
			return c.JSON(&yolov5verify.Response{
				Hint:   "未登录",
				Labels: []string{},
				Colors: []string{},
			})
		}
		var req yolov5verify.Request
		if err := c.BodyParser(&req); err != nil {
			return err
		}
		labels, colors, hint, verifyOK := yolov5verify.Verify(uid, req.Boxes, req.Scores, req.Classes)
		if verifyOK {
			flagContent, expiredAt := dynamicflag.CalcFlag(strconv.FormatUint(uid, 10), "_2024_52pojie_flagC_", time.Now())
			hint = fmt.Sprintf("flagC{%s} 过期时间: %s", flagContent, expiredAt.Format(time.DateTime))
		}
		return c.JSON(&yolov5verify.Response{
			Hint:   hint,
			Labels: labels,
			Colors: colors,
		})
	})

	// region 2048 游戏
	const Game2048CookieName = "game2048_user_data"
	game2048AuthMiddleware := func(c *fiber.Ctx) error {
		uidStr := c.Cookies("uid")
		uid, _ := strconv.ParseUint(uidStr, 10, 64)
		if uid == 0 {
			return c.JSON(map[string]any{
				"code": 1,
				"msg":  "未登录",
			})
		}
		c.Locals("uid", uid)
		return c.Next()
	}
	app.Get("/flagB/info", game2048AuthMiddleware, func(c *fiber.Ctx) error {
		return c.JSON(map[string]any{
			"code": 0,
			"msg":  "OK",
			"data": game2048.ToUserDataResponse(game2048.ParseUserData(c.Cookies(Game2048CookieName))),
		})
	})
	app.Post("/flagB/restart", game2048AuthMiddleware, func(c *fiber.Ctx) error {
		userData := game2048.ParseUserData(c.Cookies(Game2048CookieName))
		userData.NewGame()
		c.Cookie(&fiber.Cookie{
			Name:    Game2048CookieName,
			Value:   game2048.SerializeUserData(userData),
			Expires: time.Now().Add(30 * 86400 * time.Second),
		})
		return c.JSON(map[string]any{
			"code": 0,
			"msg":  "OK",
		})
	})
	app.Post("/flagB/move", game2048AuthMiddleware, func(c *fiber.Ctx) error {
		direction, _ := strconv.Atoi(c.FormValue("direction"))
		userData := game2048.ParseUserData(c.Cookies(Game2048CookieName))
		err := userData.Move(direction)
		if err != nil {
			return c.JSON(map[string]any{
				"code": 1,
				"msg":  err.Error(),
			})
		}
		c.Cookie(&fiber.Cookie{
			Name:    Game2048CookieName,
			Value:   game2048.SerializeUserData(userData),
			Expires: time.Now().Add(30 * 86400 * time.Second),
		})
		return c.JSON(map[string]any{
			"code": 0,
			"msg":  "OK",
			"data": game2048.ToUserDataResponse(userData),
		})
	})
	app.Get("/flagB/shop", game2048AuthMiddleware, func(c *fiber.Ctx) error {
		return c.JSON(map[string]any{
			"code": 0,
			"msg":  "OK",
			"data": game2048.ShopItemList,
		})
	})
	app.Post("/flagB/buy_item", game2048AuthMiddleware, func(c *fiber.Ctx) error {
		shopItemID, _ := strconv.Atoi(c.FormValue("shop_item_id"))
		buyCount, _ := strconv.ParseInt(c.FormValue("buy_count"), 10, 64)
		userData := game2048.ParseUserData(c.Cookies(Game2048CookieName))
		err := userData.BuyItem(shopItemID, buyCount)
		if err != nil {
			return c.JSON(map[string]any{
				"code": 1,
				"msg":  err.Error(),
			})
		}
		c.Cookie(&fiber.Cookie{
			Name:    Game2048CookieName,
			Value:   game2048.SerializeUserData(userData),
			Expires: time.Now().Add(30 * 86400 * time.Second),
		})
		return c.JSON(map[string]any{
			"code": 0,
			"msg":  "OK",
		})
	})
	app.Post("/flagB/use_item", game2048AuthMiddleware, func(c *fiber.Ctx) error {
		uid := c.Locals("uid").(uint64)
		itemID, _ := strconv.Atoi(c.FormValue("item_id"))
		userData := game2048.ParseUserData(c.Cookies(Game2048CookieName))
		result, err := userData.UseItem(itemID, uid)
		if err != nil {
			return c.JSON(map[string]any{
				"code": 1,
				"msg":  err.Error(),
			})
		}
		c.Cookie(&fiber.Cookie{
			Name:    Game2048CookieName,
			Value:   game2048.SerializeUserData(userData),
			Expires: time.Now().Add(30 * 86400 * time.Second),
		})
		return c.JSON(map[string]any{
			"code": 0,
			"msg":  "OK",
			"data": result,
		})
	})
	// endregion

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("X-Flag2", "flag2{xHOpRP}")
		return c.Redirect("/index.html")
	})
	app.Get("*", etag.New(etag.Config{Weak: true}), filesystem.New(filesystem.Config{
		Root:         http.Dir("./public"),
		MaxAge:       1,
		Index:        "not_exist_file",
		NotFoundFile: "",
	}))
	app.Listen(listenAddr)
}
