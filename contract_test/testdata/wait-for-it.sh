#!/bin/bash

echo "Waiting Linkd to launch on $1:$2..."

while ! nc -z $1 $2; do
  sleep 0.1 # wait for 1/10 of the second before check again
done

echo "Linkd launched"