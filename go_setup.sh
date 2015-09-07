#!/bin/bash
echo "Welcome to imgturtle setup!"
cd ..
go get github.com/gorilla/mux
go get github.com/gorilla/sessions
go get github.com/fatih/color
go get github.com/codegangsta/negroni
go get github.com/lib/pq
go get golang.org/x/crypto/pbkdf2
go get github.com/asaskevich/govalidator
echo "Finished!"
