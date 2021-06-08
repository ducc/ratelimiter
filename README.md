# Rate limiter
Limit x number of actions per minute

## Usage
This example allows 50 actions per minute and attempts to do 1000 (therefore it will take 20 minutes).
```
import "github.com/ducc/ratelimiter"

func main() {
    limiter := ratelimiter.New(50)

    for i := 0; i < 1000; i++ {
        limiter.Aquire()

        // do some work
    }
}
```