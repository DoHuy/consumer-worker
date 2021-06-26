#!/usr/bin/env bash
envsubst <configmap-$ENVIRONMENT.yml >k8s-config.yml
envsubst <main.yml >k8s-main.yml

kubectl apply -f k8s-config.yml -n ${NAMESPACE}
kubectl apply -f k8s-main.yml -n ${NAMESPACE}

kubectl rollout status deployments/${CI_PROJECT_NAME} -n ${NAMESPACE}

if [[ $? != 0 ]]; then kubectl logs -n ${NAMESPACE} -c main $(kubectl get pods -n ${NAMESPACE} --sort-by=.metadata.creationTimestamp | grep "${CI_PROJECT_NAME}" | awk '{print $1}' | tac | head -1) --tail=20 && exit 1; fi

# wait for flagger to detect the change
ok=false
until ${ok}; do
  kubectl get canary/${CI_PROJECT_NAME} -n ${NAMESPACE} | grep 'Progressing' && ok=true || ok=false
  sleep 5
done

# wait for canary analysis to finish
kubectl wait canary/${CI_PROJECT_NAME} --for=condition=promoted --timeout=15m -n ${NAMESPACE}

# check if the deployment was successful
kubectl get canary/${CI_PROJECT_NAME} -n ${NAMESPACE} | grep 'Succeeded' && exit_status=0 || exit_status=1

exit $exit_status