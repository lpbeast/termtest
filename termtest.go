package main

import (
	"fmt"
	"os"
	
	"github.com/nsf/termbox-go"
)

type Cell struct {
	Icon rune
	Walkable bool
}

func main() {
	err := termbox.Init()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer termbox.Close()
	
	var roommap [][]rune

	
	x := 2
	y := 2
	happen := true
	fg := termbox.ColorLightGray
	bg := termbox.ColorDefault
	
	roommap = append(roommap, []rune("##############"))
	roommap = append(roommap, []rune("#...#........#"))
	roommap = append(roommap, []rune("#...#....#####"))
	roommap = append(roommap, []rune("#...#........#"))
	roommap = append(roommap, []rune("#..######....#"))
	roommap = append(roommap, []rune("#............#"))
	roommap = append(roommap, []rune("####.###.##..#"))
	roommap = append(roommap, []rune("#........#...#"))
	roommap = append(roommap, []rune("#...#....#.#.#"))
	roommap = append(roommap, []rune("#...#......#.#"))
	roommap = append(roommap, []rune("#...########.#"))
	roommap = append(roommap, []rune("#............#"))
	roommap = append(roommap, []rune("##############"))
		
	rooms := make([][]Cell, len(roommap))
	
	for ry := range roommap {
		for rx := range roommap[ry] {
			var newcell Cell
			newcell.Icon = roommap[ry][rx]
			if newcell.Icon == '#' {
				newcell.Walkable = false
			} else { 
				newcell.Walkable = true
			}
			rooms[ry] = append(rooms[ry], newcell)
		}
	}
	
	for happen {
		for ry := range rooms {
			for rx := range rooms[ry] {
				termbox.SetCell(rx, ry, rooms[ry][rx].Icon, fg, bg)
			}
		}
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			switch {
			case ev.Key == termbox.KeyEsc:
				happen = false
			case ev.Key == termbox.KeyArrowUp || ev.Ch == 'w':
				if rooms[y - 1][x].Walkable { y -= 1 }
			case ev.Key == termbox.KeyArrowDown || ev.Ch == 's':
				if rooms[y + 1][x].Walkable { y += 1 }
			case ev.Key == termbox.KeyArrowLeft || ev.Ch == 'a':
				if rooms[y][x - 1].Walkable { x -= 1 }
			case ev.Key == termbox.KeyArrowRight || ev.Ch == 'd':
				if rooms[y][x + 1].Walkable { x += 1 }
			}
		}
		termbox.SetCell(x, y, '@', fg, bg)
		termbox.Flush()
	}
}