#!/bin/bash
BUILD_FOLDER=build

# create build folder
rm -rf $BUILD_FOLDER
mkdir $BUILD_FOLDER

#build driver-api
cd driver-api
go build
mv ./driver-api ../$BUILD_FOLDER
cd ..

#build driver-redis-service
cd driver-redis-service
go build
mv ./driver-redis-service ../$BUILD_FOLDER
cd ..

#build driver-store-service
cd driver-store-service
go build
mv ./driver-store-service ../$BUILD_FOLDER
cd ..


