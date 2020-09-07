# Geacon

**Using Go to implement CobaltStrike's Beacon**

----

*This project is for learning protocol analysis and reverse engineering only, if someone's rights have been violated, please contact me to remove the project, and the last DO NOT USE IT ILLEGALLY*



## How to play

1. Setup the teamserver and start a http lisenter, the teamserver will generate the file `.cobaltstrike.beacon_keys`.
2. Compile the BeaconTool with Jetbrains Idea, use command `java -jar BeaconTool.jar ` to convert java keystore to PEM format.
3. Replace the RSA key pair in the file `cmd/config/config.go` (the RSA private key is not required, I wrote it in the code just for the record)
4. Compile the geacon whatever platform you want to run: for example, use the command `export GOOS="darwin" && export GOARCH="amd64" && go build cmd/main.go` to compile an executable binary running on MacOS. 
5. Having fun ! PR and issue is welcome ;)
6. Geacon has just been tested on CobaltStrike 3.14 and only support default c2profile, so many hardcode in the project and I will not try to implement more C2profile support at this moment.
7. Thanks for **[@xxxxxyyyy](https://github.com/xxxxxyyyy)**'s PR, And now Geacon supports **CobaltStrike 4.0**, please checkout the branch `4.0` to compile.
8. Geacon's branch `master` supports **CobaltStrike 4.1**, currently available functions include: executing commands, uploading, downloading, file browser, switching the current working directory, and exiting the current process.
9. Geacon only focuses on protocol analysis, but if you want to experience more features, you can use another project of our partners, check out **[CrossC2](https://github.com/gloxec/CrossC2)** now!

## Screenshot

Get the Geacon's command execution results on Linux.
![login](https://github.com/darkr4y/geacon/raw/master/screenshots/sc.png)



## Protocol analysis

To be continued, I will update as soon as I have time ...



## Todo

1. ~~Support CobaltStrike 4.x~~

2. Fix the OS icon issue in session table
3. String encoding issue

*_DarkRay@RedCore*

