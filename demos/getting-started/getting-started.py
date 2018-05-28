#!/usr/bin/env python3

import sys
import os
import json

import dataset

def stop(msg):
    print(msg)
    sys.exit(1)

#
# create a collection with init
#
err = dataset.init("friends.ds")
if err != '':
    stop(err)
err = dataset.init("favorites.ds")
if err != '':
    stop(err)
err = dataset.init("characters.ds")
if err != '':
    stop(err)

#
# create, read, update and delete
#

## create
err = dataset.create("friends.ds", "frieda", {"name":"Little Frieda","email":"frieda@inverness.example.org"})
if err != '':
    stop(err)
err = dataset.create("friends.ds", "mojo", {"name": "Mojo Sam, the Yudoo Man", "email": "mojosam@cosmic-cafe.example.org"})
if err != '':
    stop(err)
err = dataset.create("friends.ds", "jack", {"name": "Jack Flanders", "email": "capt-jack@cosmic-voyager.example.org"})
if err != '':
    stop(err)


## read
(frieda_profile, err) = dataset.read("friends.ds", "frieda")
if err != '':
    stop(err)
(mojo_profile, err) = dataset.read("friends.ds", "mojo")
if err != '':
    stop(err)
(jack_profile, err) = dataset.read("friends.ds", "jack")
if err != '':
    stop(err)

## update

frieda_profile["catch_phrase"] = "Wowee Zowee"
mojo_profile["catch_phrase"] = "Feet Don't Fail Me Now!"
jack_profile["catch_phrase"] = "What is coming at you is coming from you"
    
err = dataset.update('friends.ds', 'frieda', frieda_profile)
if err != '':
    stop(err)
err = dataset.update('friends.ds', 'mojo', mojo_profile)
if err != '':
    stop(err)
err = dataset.update('friends.ds', 'jack', jack_profile)
if err != '':
    stop(err)

## delete

err = dataset.delete('friends.ds', 'jack')
if err != '':
    stop(err)

#
# Keys and count
#

# count
cnt = dataset.count('friends.ds')
print(f"Total Records now: {cnt}")

# keys

keys = dataset.keys('friends.ds')
print("\n".join(keys))

#
# Grids and Frames
#

# grids

keys = dataset.keys('friends.ds')
(g, err) = dataset.grid('friends.ds', keys, ['.name', '.email', 'catch_phrase'])
if err != '':
    stop(err)
print(json.dumps(g, indent = 4))

# frames


