#!/usr/bin/python

from subprocess import call
import time, sys


args = ["wbm", "10","p"]
for n in range(1,11):
	args[1] = str(n * 100)
	start_time = time.time()
	call(args)
	print "execution took ",time.time() - start_time, "seconds"
