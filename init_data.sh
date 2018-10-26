#!/bin/bash

mongo reunion --eval "db.announcement.remove({})"

indexGlobal=1
for j in {1..100}
do
	insertValue=""
	for i in {1..10}
	do
		insertValue=$insertValue"{'name': 'Médor $i$j', 'description' : 'description test','race': 'Shiba Inu', 'city': 'Saint Denis', 'color': 'Noir et blanc', 'phonenumber' : '0606060606','account': {'email': 'test@test', 'password': 'totoPassword'}, 'type' : 'lost', 'state': 'validated', 'creationdate' : ISODate('20$i$j-02-21T15:00:00Z')},"	
		indexGlobal=$indexGlobal+1
	done
	insertValue="db.announcement.insert([$insertValue])"
	
	echo $insertValue | mongo reunion
done
for j in {1..3}
do
	insertValue=""
	for i in {1..10}
	do
		insertValue=$insertValue"{'name': 'Toto $i$j', 'description' : 'description test','race': 'Shiba Inu', 'city': 'Saint Pierre', 'color': 'Noir et blanc', 'phonenumber' : '0606060606','account': {'email': 'test@test', 'password': 'totoPassword'}, 'type' : 'lost', 'state': 'validated', 'creationdate' : ISODate('20$i$j-02-21T15:00:00Z')},"	
		indexGlobal=$indexGlobal+1
	done
	insertValue="db.announcement.insert([$insertValue])"
	
	echo $insertValue | mongo reunion
done
