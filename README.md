# dingrobot
DingTake Robot API

Prepare DingTalk robot access_token, See my article <https://testerhome.com/topics/11217> if you dont know how to get it.

## 注意

当前使用keyword 和 Secret 两种方式

两种方式的确定，决定于创建钉钉机器人时使用了何种验证

**使用**
1. keyword
>robot.SetKeyWord("111")

2. Secret
>robotSecret.SetSecret("SEC4dd5f7859340f03f9ce3db9affd27ac269b991f2a7b81b7cf341c4de6173024e")


## Usage
```go
import "github.com/mostcute/dingrobot"

func main(){
//测试key word
robot := dingrobot.New("fa1e90877ab61046ef95c357f3660faa80e1dbc4e50aa108c1da88539a544392")
robot.SetKeyWord("111")

// send text message
robot.Text("test keyword")
// send markdown message
robot.Markdown("makedown","**test keyword**")
//send link
robot.Link("Google", "Google homepage", "https://www.google.com.hk","https://www.google.com.hk")

// At someone
robot.At("18658148376").Text("test keyword")

// At all
robot.AtAll(true).Text("test keyword")



//测试 Secret
robotSecret := dingrobot.New("c1a81c5a7ff1411ef5f6c074cbc8d71c464bc4ccc7eb86bc25ba3596f7ec703b")
robotSecret.SetSecret("SEC4dd5f7859340f03f9ce3db9affd27ac269b991f2a7b81b7cf341c4de6173024e")
// send text message
robotSecret.Text("robotSecret")
// send markdown message
robotSecret.Markdown("makedown","**robotSecret**")
// send link
robotSecret.Link("Google", "Google homepage", "https://www.google.com.hk","https://www.google.com.hk")

// At someone
robotSecret.At("18658148376").Text("robotSecret")

// At all
robotSecret.AtAll(true).Text("robotSecret")
}
```

## TODO
* FeedCard
* ActionCard
* ImageUpload

## LICENSE
[MIT](LICENSE)