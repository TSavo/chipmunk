package main

import (
	"encoding/json"
	//"fmt"
	"github.com/TSavo/chipmunk"
	"github.com/TSavo/chipmunk/vect"
	"math"
	"math/rand"
	"net"
	//"time"
)

type Ship struct {
	Id    int
	Shape *chipmunk.Shape
}

var (
	ballRadius  = 25
	ballMass    = 1
	space       *chipmunk.Space
	balls       []*Ship
	staticLines []*chipmunk.Shape
	deg2rad     = math.Pi / 180
)

func addBall() {
	x := rand.Intn(1135) + 115
	//ball := chipmunk.NewCircle(vect.Vector_Zero, vect.Float(ballRadius))
	ball := chipmunk.NewBox(vect.Vector_Zero, 50, 50)
	ball.SetElasticity(0.9)
	body := chipmunk.NewBody(vect.Float(10000000), ball.Moment(vect.Float(100000000)))
	body.SetPosition(vect.Vect{vect.Float(x), 100.0})
	body.SetAngle(vect.Float(rand.Float64() * 2 * math.Pi))
	//	t := 2
	//	if(rand.Intn(2) == 1){
	//		t *= -1
	//	}
	//	body.AddTorque(vect.Float(t))
	body.AddShape(ball)
	space.AddBody(body)
	balls = append(balls, &Ship{x, ball})
}

// step advances the physics engine and cleans up any balls that are off-screen
func step(dt float32) {
	space.Step(vect.Float(dt))
	for i := 0; i < len(balls); i++ {
		p := balls[i].Shape.Body.Position()
		if p.Y > 1500 {
			space.RemoveBody(balls[i].Shape.Body)
			balls = append(balls[:i], balls[i+1:]...)
			i-- // consider same index again
		}
	}

}

// createBodies sets up the chipmunk space and static bodies
func createBodies() {
	space = chipmunk.NewSpace()
	space.Gravity = vect.Vect{0, 2900}

	staticBody := chipmunk.NewBodyStatic()
	staticLines = []*chipmunk.Shape{
		chipmunk.NewSegment(vect.Vect{100.0, 500.0}, vect.Vect{1000.0, 500.0}, 0),
		chipmunk.NewSegment(vect.Vect{1000.0, 500.0}, vect.Vect{1000.0, 100.0}, 0),
	}
	for _, segment := range staticLines {
		segment.SetElasticity(0.9)
		staticBody.AddShape(segment)
	}
	space.AddBody(staticBody)
}

type Box struct {
	Id int
	X, Y, A vect.Float
}

type Line struct {
	Point1, Point2 vect.Vect
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		encoder := json.NewEncoder(conn)
		createBodies()
		ticksToNextBall := 10
		//ticker := time.NewTicker(time.Second / 60)
		for {
			ticksToNextBall--
			if ticksToNextBall == 0 {
				ticksToNextBall = rand.Intn(20) + 1
				addBall()
			}
			step(1.0 / 60.0)
			
			message := make(map[string]interface{})
			
			boxes := make([]Box, len(balls))
			
			for i, x := range balls {
				boxes[i] = Box{x.Id, x.Shape.Body.Position().X, x.Shape.Body.Position().Y, x.Shape.Body.Angle()}
			}
			lines := make([]Line, len(staticLines))
			for i, x := range staticLines {
				lines[i].Point1 = x.GetAsSegment().A
				lines[i].Point2 = x.GetAsSegment().B
			}
			message["boxes"] = boxes
			message["lines"] = lines
			encoder.Encode(message)
			//<-ticker.C // wait up to 1/60th of a second
		}
	}
	// set up physics
}
