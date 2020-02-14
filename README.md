# connect-dots

The goal of the game is to cover the board with segments connecting the dots.
Two dots having the same color can be connected with a segment.
The segments should not intersects each other.

The application is not (and will not be) a commercial grade 'game' (not a casual game either) but a simple proof of concept which allowed me to check how fast or easy I could implement a simple board game using golang (https://golang.org/) and libsdl (https://www.libsdl.org/). 

The Go binding for libsdl may be found at this location: https://github.com/veandco/go-sdl2 .

The game could be built statically against go-sdl2 libraries by running the followin command (here the OS is linux): 

CGO_ENABLED=1 CC=gcc GOOS=linux GOARCH=amd64 go build -tags static -ldflags "-s -w"

# Screenshoots
![connect-dots-3](https://user-images.githubusercontent.com/59707990/74368823-09751280-4ddd-11ea-9c28-47c72c4d2814.png)
![connect-dots-5](https://user-images.githubusercontent.com/59707990/74548280-3862c400-4f56-11ea-85c8-20ee09586ad8.png)

