#!/bin/bash

unset KUBECONFIG
CLUSTERNAME=`cat kind.yaml| grep "name" | awk '{print $2}'`
kind get clusters | grep "$CLUSTERNAME"

if [ "$?" != "0" ]; then
	echo -e "cluster $CLUSTERNAME not found|\ncreating ..."
	kind create cluster --config kind.yaml

	echo "waiting the cluster started ..."
	attemptsMax=5
	attempts=0

	while [ $attemptsMax -gt $attempts ]
	do
		kubectl get nodes -o json | jq '.items[].status[][]' | grep '\"status\": \"True\"'
		if [ "$?" -eq "0" ]; then
			echo "Cluster Ready"
			break
			exit 0
		fi
		sleep 3
		attemps=$((attemps + 1))
	done
	exit 1

fi

exit 0
