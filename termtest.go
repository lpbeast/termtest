package main

import (
	"fmt"
	"os"
	"github.com/nsf/termbox-go"
)

type Tile struct {
	Icon rune
	Walkable bool
	Pit bool
}

type Block struct {
	X, Y int
}

type Player struct {
	X, Y, Zone int
}

func Tbprintf(x, y int, fg, bg termbox.Attribute, s string, a ...any) {
	compos := fmt.Sprintf(s, a...)
	for i := 0; i < len(compos); i++{
		r := rune(compos[i])
		termbox.SetCell(x + i, y, r, fg, bg)
	}
}

func (p *Player) Move(dir rune) {
	switch {
	case dir == 'n':
		p.Y -= 1
	case dir == 's':
		p.Y += 1
	case dir == 'w':
		p.X -= 1
	case dir == 'e':
		p.X += 1
	default:
	}
}

func Boxcheck(x, y int, dir rune, b *[]Block, r [][]Tile) bool {
	var tx, ty, nx, ny int
	canmove := true
	switch {
	case dir == 'n':
		tx = x
		ty = y - 1
		nx = x
		ny = y - 2
	case dir == 's':
		tx = x
		ty = y + 1
		nx = x
		ny = y + 2
	case dir == 'e':
		tx = x + 1
		ty = y
		nx = x + 2
		ny = y
	case dir == 'w':
		tx = x - 1
		ty = y
		nx = x - 2
		ny = y
	default:
	}
	for k, i := range(*b) {
		if i.X == tx && i.Y == ty { //there's a box there
			fell := false
			for _, j := range(*b) {
				if j.X == nx && j.Y == ny { //there's a box behind the box
					canmove = false
				}
				if !r[ny][nx].Walkable {
					if !r[ny][nx].Pit {
						canmove = false
					} else {
						*b = append((*b)[:k], (*b)[k + 1:]...)
						r[ny][nx].Pit = false
						r[ny][nx].Walkable = true
						r[ny][nx].Icon = '.'
						fell = true
					}
				}
			}
			if canmove && !fell{
				(*b)[k].X = nx
				(*b)[k].Y = ny
			}
		}
	}
	return canmove
}

func (p *Player) Moveto(x, y, z int) {
	p.X = x
	p.Y = y
	p.Zone = z
}

func main() {
	err := termbox.Init()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer termbox.Close()
	
	var roommap [][]rune
	var yourchar Player
	yourchar.X = 2
	yourchar.Y = 2
	yourchar.Zone = 1
	happen := true
	fg := termbox.ColorLightGray
	bg := termbox.ColorDefault
	stepct := 0
	
	roommap = append(roommap, []rune("##############"))
	roommap = append(roommap, []rune("#...#........#"))
	roommap = append(roommap, []rune("#...#.O_..####"))
	roommap = append(roommap, []rune("#...#........#"))
	roommap = append(roommap, []rune("#..#####..O..#"))
	roommap = append(roommap, []rune("#............#"))
	roommap = append(roommap, []rune("####.###O##..#"))
	roommap = append(roommap, []rune("#........#...#"))
	roommap = append(roommap, []rune("#...#....#.#.#"))
	roommap = append(roommap, []rune("#...#......#.#"))
	roommap = append(roommap, []rune("#...########.#"))
	roommap = append(roommap, []rune("#............#"))
	roommap = append(roommap, []rune("##############"))
		
	rooms := make([][]Tile, len(roommap))
	boxes := make([]Block, 0)
	
	for ry := range roommap {
		for rx := range roommap[ry] {
			var newtile Tile
			newtile.Pit = false
			newtile.Icon = roommap[ry][rx]
			if newtile.Icon == '#' {
				newtile.Walkable = false
			} else { 
				if newtile.Icon == '_' {
					newtile.Walkable = false
					newtile.Pit = true
				} else {
					newtile.Walkable = true
					if newtile.Icon == 'O' {
						newtile.Icon = '.'
						var b Block
						b.X = rx
						b.Y = ry
						boxes = append(boxes, b)
					}
				}
			}
			rooms[ry] = append(rooms[ry], newtile)
		}
	}

	for happen {
		for ry := range rooms {
			for rx := range rooms[ry] {
				termbox.SetCell(rx, ry, rooms[ry][rx].Icon, fg, bg)
			}
		}
		termbox.SetCell(yourchar.X, yourchar.Y, '@', fg, bg)
		for _, i := range boxes {
			termbox.SetCell(i.X, i.Y, 'O', fg, bg)
		}
		Tbprintf(0, 14, fg, bg, "Number of steps taken: %d", stepct)
		termbox.Flush()
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			stepct++
			x := yourchar.X
			y := yourchar.Y
			switch {
			case ev.Key == termbox.KeyEsc:
				happen = false
			case ev.Ch == ' ' || ev.Ch == '.':
			case ev.Key == termbox.KeyArrowUp || ev.Ch == 'w':
				if rooms[y - 1][x].Walkable { 
					if Boxcheck(x, y, 'n', &boxes, rooms) {
						yourchar.Move('n') 
					}
				}
			case ev.Key == termbox.KeyArrowDown || ev.Ch == 's':
				if rooms[y + 1][x].Walkable {
					if Boxcheck(x, y, 's', &boxes, rooms) {
						yourchar.Move('s') 
					}
				}		
			case ev.Key == termbox.KeyArrowLeft || ev.Ch == 'a':
				if rooms[y][x - 1].Walkable { 
					if Boxcheck(x, y, 'w', &boxes, rooms) {
						yourchar.Move('w') 
					}
				}
			case ev.Key == termbox.KeyArrowRight || ev.Ch == 'd':
				if rooms[y][x + 1].Walkable { 
					if Boxcheck(x, y, 'e', &boxes, rooms) {
						yourchar.Move('e') 
					}
				}
			}
		}
	}
}