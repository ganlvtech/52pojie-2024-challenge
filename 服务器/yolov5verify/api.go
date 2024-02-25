package yolov5verify

var Labels = []string{
	"person",
	"bicycle",
	"car",
	"motorcycle",
	"airplane",
	"bus",
	"train",
	"truck",
	"boat",
	"traffic light",
	"fire hydrant",
	"stop sign",
	"parking meter",
	"bench",
	"bird",
	"cat",
	"dog",
	"horse",
	"sheep",
	"cow",
	"elephant",
	"bear",
	"zebra",
	"giraffe",
	"backpack",
	"umbrella",
	"handbag",
	"tie",
	"suitcase",
	"frisbee",
	"skis",
	"snowboard",
	"sports ball",
	"kite",
	"baseball bat",
	"baseball glove",
	"skateboard",
	"surfboard",
	"tennis racket",
	"bottle",
	"wine glass",
	"cup",
	"fork",
	"knife",
	"spoon",
	"bowl",
	"banana",
	"apple",
	"sandwich",
	"orange",
	"broccoli",
	"carrot",
	"hot dog",
	"pizza",
	"donut",
	"cake",
	"chair",
	"couch",
	"potted plant",
	"bed",
	"dining table",
	"toilet",
	"tv",
	"laptop",
	"mouse",
	"remote",
	"keyboard",
	"cell phone",
	"microwave",
	"oven",
	"toaster",
	"sink",
	"refrigerator",
	"book",
	"clock",
	"vase",
	"scissors",
	"teddy bear",
	"hair drier",
	"toothbrush",
}

type Request struct {
	Boxes   []float64 `json:"boxes"`
	Scores  []float64 `json:"scores"`
	Classes []int     `json:"classes"`
}

type Response struct {
	Hint   string   `json:"hint"`
	Labels []string `json:"labels"`
	Colors []string `json:"colors"`
}

type Item struct {
	Index       int
	X1          float64
	Y1          float64
	X2          float64
	Y2          float64
	Score       float64
	Class       int
	ClassExists bool
	PositionOK  bool
	Label       string
	Color       string
}

func CalcAnswer(uid uint64, answerLen int) []int {
	var result []int
	seed := uid & 0x7fffffff
	for i := 0; i < answerLen; i++ {
		seed = (seed*1103515245 + 12345) & 0xffffffff
		result = append(result, int(seed%10))
	}
	return result
}

func GetItemsBoundingBox(items []*Item) (minX float64, minY float64, maxX float64, maxY float64) {
	minX = 1
	minY = 1
	maxX = 0
	maxY = 0
	for _, item := range items {
		if item.X1 < minX {
			minX = item.X1
		}
		if item.X2 < minX {
			minX = item.X2
		}
		if item.X1 > maxX {
			maxX = item.X1
		}
		if item.X2 > maxX {
			maxX = item.X2
		}
		if item.Y1 < minY {
			minY = item.Y1
		}
		if item.Y2 < minY {
			minY = item.Y2
		}
		if item.Y1 > maxY {
			maxY = item.Y1
		}
		if item.Y2 > maxY {
			maxY = item.Y2
		}
	}
	return
}

func GetItemPosition(x1 float64, y1 float64, x2 float64, y2 float64) int {
	x := (x1 + x2) / 2.0
	y := (y1 + y2) / 2.0
	if x < 0.5 && y < 0.5 {
		return 0
	}
	if x >= 0.5 && y < 0.5 {
		return 1
	}
	if x < 0.5 && y >= 0.5 {
		return 2
	}
	if x >= 0.5 && y >= 0.5 {
		return 3
	}
	return -1
}

func VerifyItems(uid uint64, items []*Item) (hint string, verifyOK bool) {
	// 检查物体总数量
	if len(items) > 4 {
		hint = "物体太多了"
	} else if len(items) < 4 {
		hint = "物体太少了"
	}
	// 检查物体种类
	answer := CalcAnswer(uid, 4)
	answerMap := map[int]map[int]int{} // map[class]map[position]count
	for i, v := range answer {
		if _, ok := answerMap[v]; !ok {
			answerMap[v] = map[int]int{}
		}
		answerMap[v][i] += 1
	}
	minX, minY, maxX, maxY := GetItemsBoundingBox(items)
	dx := maxX - minX
	dy := maxY - minY
	for _, item := range items {
		if positionMap, ok := answerMap[item.Class]; ok {
			item.ClassExists = true
			if dx > 0 && dy > 0 {
				itemPosition := GetItemPosition((item.X1-minX)/dx, (item.Y1-minY)/dy, (item.X2-minX)/dx, (item.Y2-minY)/dy)
				if itemPosition >= 0 {
					// 检查物体位置
					if count, ok := positionMap[itemPosition]; ok {
						if count > 0 {
							item.PositionOK = true
							positionMap[itemPosition] -= 1
						}
					}
				}
			}
		}
	}
	for _, item := range items {
		if item.ClassExists {
			if item.PositionOK {
				item.Label = item.Label + " 种类正确 位置正确"
				item.Color = "99ff99"
			} else {
				item.Label = item.Label + " 种类正确 位置错误"
				item.Color = "ffff99"
			}
		} else {
			item.Label = item.Label + " 种类错误"
			item.Color = "ff9999"
		}
	}
	if len(items) == 4 {
		verifyOK = true
		for _, item := range items {
			if !item.PositionOK {
				verifyOK = false
				break
			}
		}
		if !verifyOK {
			hint = "数量正确 种类或位置仍需调整"
		}
	}
	return
}

func Verify(uid uint64, boxes []float64, scores []float64, classes []int) (labels []string, colors []string, hint string, verifyOK bool) {
	l := len(scores)
	if len(boxes)/4 < l {
		l = len(boxes) / 4
	}
	if len(classes) < l {
		l = len(classes)
	}
	if l > 100 {
		l = 100
	}

	var items []*Item
	for i := 0; i < l; i++ {
		if scores[i] > 0.25 && scores[i] <= 1.0 &&
			boxes[4*i+0] >= 0.0 && boxes[4*i+0] <= 1.0 &&
			boxes[4*i+1] >= 0.0 && boxes[4*i+1] <= 1.0 &&
			boxes[4*i+2] >= 0.0 && boxes[4*i+2] <= 1.0 &&
			boxes[4*i+3] >= 0.0 && boxes[4*i+3] <= 1.0 &&
			classes[i] >= 0 && classes[i] < len(Labels) {
			items = append(items, &Item{
				Index:       i,
				X1:          boxes[4*i+0],
				Y1:          boxes[4*i+1],
				X2:          boxes[4*i+2],
				Y2:          boxes[4*i+3],
				Score:       scores[i],
				Class:       classes[i],
				ClassExists: false,
				PositionOK:  false,
				Label:       Labels[classes[i]],
				Color:       "",
			})
		}
	}

	hint, verifyOK = VerifyItems(uid, items)
	labels = make([]string, len(scores))
	colors = make([]string, len(scores))
	for _, item := range items {
		labels[item.Index] = item.Label
		colors[item.Index] = item.Color
	}
	return
}
