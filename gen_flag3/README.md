# 噪声图片相对运动

```bash
go run main.go
cd output
ffmpeg -r 60 -f image2 -i flag3_%02d.png flag3.mp4
```
