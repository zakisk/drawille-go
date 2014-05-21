package main

import (
  "github.com/exrook/drawille-go"
  "math"
  "fmt"
)

func main() {
  s := drawille.NewCanvas()
  for x:=0;x<(1800);x=x+1 {
    y := int(math.Sin((math.Pi/180)*float64(x))*10)
    s.Set(x/10,y)
  }
  fmt.Print(s.Frame(s.MinX()*2,s.MinY()*4,s.MaxX()*2,s.MaxY()*4))
  
  s.Clear()
  
  for x:=0;x<1800;x=x+10 {
    s.Set(x/10,int(10+math.Sin((math.Pi/180)*float64(x))*10))
    s.Set(x/10,int(10+math.Cos((math.Pi/180)*float64(x))*10))
  }
  fmt.Print(s.Frame(s.MinX()*2,s.MinY()*4,s.MaxX()*2,s.MaxY()*4))
  
  s.Clear()
  
  for x:=0;x<3600;x=x+20 {
    s.Set(x/20,int(4+math.Sin((math.Pi/180)*float64(x))*4))
  }  
  fmt.Print(s.Frame(s.MinX()*2,s.MinY()*4,s.MaxX()*2,s.MaxY()*4))
  
  s.Clear()

  for x:=0;x<360;x=x+1 {
    s.Set(x/4,int(30+math.Sin((math.Pi/180)*float64(x))*30))
  }
  
  for x:=0;x<30;x=x+1 {
    for y:=0;y<30;y=y+1 {
      s.Set(x,y)
      s.Toggle(x+30,y+30)
      s.Toggle(x+60,y)
    }
  }
  fmt.Print(s.Frame(s.MinX()*2,s.MinY()*4,s.MaxX()*2,s.MaxY()*4))  
}
