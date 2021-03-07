## Lattice Boltzman Solver

This was a project that was created primarily to learn golang. It uses [gomobile](https://github.com/golang/mobile) to build artifacts for android and ios, the application can also run on mac, linux and windows.

**NOTE**: This project is unmaintained - I have uploaded the code for reference.

### Usage
  - Install pre-requisites
    + Install golang
    + Install gomobile `go get golang.org/x/mobile/cmd/gomobile`
    + Initialize gomobile `gomobile init` (This could take a few minutes)
  - Build on windows, linux or mac
    + `go build .`
    + Run the `go_lbm` binary
  - Build mobile applications
    + iOS - `gomobile build -target=ios .`
    + Android - `gomobile build -target=android .`
    + For more information on building for mobile applications, refer to gomobile docs.
  - Alternatively, use `make` along with the provided `Makefile`

### Demo
![Go LBM Demo](https://github.com/prasadchandan/go_lbm/blob/master/repo-assets/demo.gif)

### Credits

 - The code for the solver was ported based on the JavaScript implementation by Dan Schroeder 
   + https://physics.weber.edu/schroeder/fluids/
 - Carlo Barth for the colorMapCreator python script that was used to generate the color maps used in the application
 - [FiraSans Font](https://github.com/mozilla/Fira)
