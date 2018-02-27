package compiler

// chan_test.go: tests of the all-lua channels

import (
	"fmt"
	"testing"
	"time"

	//"github.com/gijit/gi/pkg/token"
	//"github.com/gijit/gi/pkg/types"
	cv "github.com/glycerine/goconvey/convey"
	//"github.com/glycerine/luar"
)

func Test900SendAndRecvAllLu(t *testing.T) {

	cv.Convey("select{} should block the goroutine forever, unless also select{default:}", t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		// with default: present we should not block
		// _selection = __task.select({{}});
		code := ` a:= 0; go func() { println("top of go-started func"); a = 1; select{ default: }; a= 2; }() // should not block`
		translation, err := inc.Tr([]byte(code))
		//*dbg = true
		fmt.Printf("translation='%s'\n", string(translation))

		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "a", 2)

		//  _r = __task.select({});
		code = ` b:= 0; go func() { b = 1; select{}; b= 2; }() // should block goroutine forever`
		translation, err = inc.Tr([]byte(code))
		panicOn(err)
		fmt.Printf("translation='%s'\n", string(translation))

		LuaRunAndReport(vm, string(translation))
		select {
		case <-time.After(1 * time.Second):
		}
		LuaMustInt64(vm, "b", 1)
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test901(t *testing.T) {

	cv.Convey("In the all-lua go/coroutine system, ch := make(chan int, 1); ch <- 56;  b := <-ch; write and read back b of 57", t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		// with default: present we should not block
		// _selection = __task.select({{}});
		code := ` ch := make(chan int, 1); ch <- 56;  b := <-ch; `
		translation, err := inc.Tr([]byte(code))
		//*dbg = true
		fmt.Printf("translation='%s'\n", string(translation))

		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "b", 56)
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test902(t *testing.T) {

	cv.Convey("spawn goroutine, send and receive on unbuffered channel, in the all-lua system.", t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		code := ` ch := make(chan int); go func() {ch <- 56;}(); b := <-ch; `
		translation, err := inc.Tr([]byte(code))
		//*dbg = true
		fmt.Printf("translation='%s'\n", string(translation))

		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "b", 56)
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test903(t *testing.T) {

	cv.Convey("select on multiple channels, in the all-lua system.", t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		code := `
    a := 0
    b := ""

    chInt := make(chan int)
    chStr := make(chan string)

	go func() {
		chInt <- 43
	}()
	go func() {
		chStr <- "hello select"
	}()

    for i := 0; i < 2; i++ {
      select {
        case a = <- chInt:
        case b = <- chStr:
      }
    }
`

		translation, err := inc.Tr([]byte(code))
		//*dbg = true
		fmt.Printf("translation='%s'\n", string(translation))

		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "a", 43)
		LuaMustString(vm, "b", "hello select")
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test904(t *testing.T) {

	cv.Convey("select on with both receive and send, in the all-lua system.", t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		code := `
    a := 0
    b := ""
    sentAndReceived := ""
    chInt := make(chan int)
    chStr := make(chan string)
    chStr2 := make(chan string)

	go func() {
		chInt <- 43
	}()
	go func() {
		chStr <- "hello select"
	}()
	go func() {
		sentAndReceived = <-chStr2
	}()

    for i := 0; i < 3; i++ {
      select {
        case a = <- chInt:
        case b = <- chStr:
        case chStr2 <- "yeehaw":
      }
    }
`

		translation, err := inc.Tr([]byte(code))
		//*dbg = true
		fmt.Printf("translation='%s'\n", string(translation))

		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "a", 43)
		LuaMustString(vm, "b", "hello select")
		LuaMustString(vm, "sentAndReceived", "yeehaw")
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test905(t *testing.T) {

	cv.Convey("simple send, in the all-lua system.", t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		code := `
    sentAndReceived := ""
    chStr2 := make(chan string)

	go func() {
		sentAndReceived = <-chStr2
	}()
    chStr2 <- "yeehaw"
`

		translation, err := inc.Tr([]byte(code))
		//*dbg = true
		fmt.Printf("translation='%s'\n", string(translation))

		LuaRunAndReport(vm, string(translation))

		LuaMustString(vm, "sentAndReceived", "yeehaw")
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test906(t *testing.T) {

	cv.Convey("all-lua system: multiple sends from new goroutine, received on main goroutine", t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		code := `
    ch := make(chan int)
    start := 7
    stop  := 9
	go func() {
		for i:=start;i<stop;i++ {
           println("before sending i=", i);
           ch <- i;
           println("after sending i=", i);
        }
    }()

    a := <- ch
    b := <- ch
`

		translation, err := inc.Tr([]byte(code))
		//*dbg = true
		fmt.Printf("translation='%s'\n", string(translation))

		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "a", 7)
		LuaMustInt64(vm, "b", 8)
		//LuaMustInt64(vm, "c", 9)
		//LuaMustInt64(vm, "d", 10)
		//LuaMustInt64(vm, "e", 11)
		cv.So(true, cv.ShouldBeTrue)
	})
}