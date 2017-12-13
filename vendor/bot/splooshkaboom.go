package bot

import (
  "math/rand"
  "fmt"
  "time"
)

const (
  SKEMPTY = 0
  SKMIDDLEH = 1
  SKMIDDLEV = 8
  SKBACKW = 2
  SKBACKE = 3
  SKBACKS = 5
  SKBACKN = 6
  SKFRONTW = 7
  SKFRONTE = 9
  SKFRONTS = 10
  SKFRONTN = 11

  SKSPLOOSH = 12
  SKKABOOM = 13

  SKBOAT = 14 // this is a boat before rendering
)

const (
  BOATS = 0
  BOATW = 1
  BOATE = 2
  BOATN = 3
)

const (
  MAPW = 8
  MAPH = 10
)

const (
  BOATAMOUNT = 3
)

const MAXFAILEDVALID = 10

type BoatTile struct {
  X int
  Y int

  Value int
}

type Boat struct {
  BoatTile []BoatTile
  Orientation int
  BoatSize int
}

/*
SplooshKaboom - holds the SplooshKaboom game
GameField - 10 by 8 playing field of ints
1 is a ship, 2 is ships rear end 3 is ships front
Bombs - amount of bombs in a game 3 by 8 bombs
*/
type SplooshKaboom struct {
  GameField [MAPW][MAPH]int
  Bombs int
  BoatsLeft int
  ID string

  Boats [BOATAMOUNT]Boat
}

func NewSplooshKaboom(gid string) *SplooshKaboom {
  var newSK *SplooshKaboom = new(SplooshKaboom)

  newSK.ID = gid

  newSK.GenerateNewGame()

  return newSK
}


func (sk *SplooshKaboom) GenerateNewGame() {
  // seed
  seed := int64(time.Now().Unix())
  fmt.Printf("Seeding rand with %d\n", seed)
  rand.Seed(seed)

  for x, _ := range sk.GameField {
    for y, _ := range sk.GameField[x] {
      sk.GameField[x][y] = 0
    }
  }

  sk.BoatsLeft = BOATAMOUNT

  for i, _ := range sk.Boats {
    sk.Boats[i].Orientation = rand.Intn(4) // orientation

    var boatSize = 0
    if i == 0 {
      boatSize = 2
    } else if i == 1 {
      boatSize = 4
    } else if i == 2 {
      boatSize = 5
    }

    sk.Boats[i].BoatSize = boatSize

    sk.appendBoat(&sk.Boats[i], boatSize)
  }

  sk.Bombs = 8 * 3
}

func (sk *SplooshKaboom) appendBoat(boat *Boat, boatSize int) {
  var newX = -1
  var newY = -1

  var failedValid = 0

  // generate origin tiles (boat expands away from this tile in orientation)
  // do this until a good tile is reached that does not intersect anything
  for j := 0; j < boatSize; j++ {
    if j == 0 {
      boat.BoatTile = make([]BoatTile, 0)
      newX = rand.Intn(MAPW)
      newY = rand.Intn(MAPH)
    }
    var newBT BoatTile
    if boat.Orientation == BOATN {
      // upwards expansion
      newBT = BoatTile{X: newX - j, Y: newY}

      if j == 0 {
        newBT.Value = SKBACKN
      } else if j == boatSize - 1 {
        newBT.Value = SKFRONTN
      } else {
        newBT.Value = SKMIDDLEV
      }
    } else if boat.Orientation == BOATS {
      // upwards expansion
      newBT = BoatTile{X: newX + j, Y: newY}

      if j == 0 {
        newBT.Value = SKBACKS
      } else if j == boatSize - 1 {
        newBT.Value = SKFRONTS
      } else {
        newBT.Value = SKMIDDLEV
      }
    } else if boat.Orientation == BOATW {
      // upwards expansion
      newBT = BoatTile{X: newX, Y: newY - j}

      if j == 0 {
        newBT.Value = SKBACKW
      } else if j == boatSize - 1 {
        newBT.Value = SKFRONTW
      } else {
        newBT.Value = SKMIDDLEH
      }
    } else if boat.Orientation == BOATE {
      // upwards expansion
      newBT = BoatTile{X: newX, Y: newY + j}

      if j == 0 {
        newBT.Value = SKBACKE
      } else if j == boatSize - 1 {
        newBT.Value = SKFRONTE
      } else {
        newBT.Value = SKMIDDLEH
      }
    } else {
      fmt.Println("Warning: Unknown orientation for boat: ", newX, newY, boat.Orientation)
      newBT = BoatTile{X: newX, Y: newY}
      newBT.Value = SKMIDDLEH
    }

    if failedValid > MAXFAILEDVALID {
      fmt.Println("Failed validition too many times. Gave up.")
      break // just give up
    }

    if !sk.isTileValid(&newBT) {
      // try again
      j = 0
      // and return, because this call is bad
      failedValid++
      continue
    } else {
      boat.BoatTile = append(boat.BoatTile, newBT)
    }
  }

  boat.BoatSize = len(boat.BoatTile)
}

func (sk *SplooshKaboom) isTileValid(newBT *BoatTile) bool {
  for _, oboat := range sk.Boats {
    for _, tile := range oboat.BoatTile {
      if (newBT.X == tile.X && newBT.Y == tile.X) || newBT.X < 0 || newBT.X >= MAPW ||
      newBT.Y < 0 || newBT.Y >= MAPH {
        // if this is the case, try again!
        return false
      }
    }
    // we are good to go!
  }

  return true
}

func (sk *SplooshKaboom) IsFieldEmpty(x int, y int) (bool, int) {
  if len(sk.GameField) <= x {
    return false, -1
  }

  if len(sk.GameField[x]) <= y {
    return false, -1
  }

  // check for boats
  for _, oboat := range sk.Boats {
    for _, tile := range oboat.BoatTile {
      if tile.X == x && tile.Y == y {
        return false, sk.GameField[x][y]
      }
    }
    // we are good to go!
  }

  if sk.GameField[x][y] == SKEMPTY {
    return true, sk.GameField[x][y]
  }

  return false, sk.GameField[x][y]
}

func (sk *SplooshKaboom) Target(x int, y int) {
  empty, val := sk.IsFieldEmpty(x, y)
  if val >= 0 {
    sk.Bombs--
    if val == SKEMPTY || empty {
      sk.GameField[x][y] = SKSPLOOSH
    } else {
      sk.GameField[x][y] = SKKABOOM
    }
  }

  // check boats

  for index, oboat := range sk.Boats {
    for _, tile := range oboat.BoatTile {
      if sk.GameField[tile.X][tile.Y] == SKKABOOM {
        oboat.BoatSize--
        if oboat.BoatSize == 0 {
          sk.BoatsLeft--
          oboat.BoatSize--
          sk.Boats[index] = oboat
          fmt.Println("Boatsize ", oboat.BoatSize, oboat.BoatSize == 0, oboat, " Remaining boats: ", sk.BoatsLeft)
        }
      }
    }
  }
}

func (sk *SplooshKaboom) GameOver() {
  for _, oboat := range sk.Boats {
    for _, tile := range oboat.BoatTile {
      if sk.GameField[tile.X][tile.Y] != SKKABOOM {
        sk.GameField[tile.X][tile.Y] = tile.Value
      }
    }
  }

  fmt.Println("Revealed: ", sk.GameField)
}
