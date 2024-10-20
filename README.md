# Sound Position Timer

## Usage
* Use this tool to measure the duration of sound positions within syllables. Measure sounds for 2 second, 1 second, and half second syllable durations. Tools of this variety are frequently used in speech Pathology for patients with a stutter. 

## Setup 
* Download [Golang](https://go.dev)
* Clone this repository
```bash
$ git clone https://github.com/IMSterling/speechTimer.git
```
* Build the main file
```bash
$ go build -o speechTimer main.go
```
* Run the built file
```bash
$ ./speechTimer
``` 

Optional
* Add the directory to your system's PATH so you can run the executable from anywhere

## Implementation details

* This tool is build using [Gio](https://gioui.org/) and is designed to be platform agnostic. It runs on MacOSX, Linux, and Windows. 
* This implementation borrows heavily from [this](https://jonegil.github.io/gui-with-gio/egg_timer/) fantastic Gio tutorial 

## Contact 
* Got questions? Email me at ianmcdiarmidsterling at gmail dot com
