Hello team!

I had a lot of fun with this exercise and some unexpected challenges, but I'm looking forward to discussing this experience further with you!

There are two directories in this repo. There is a k8s directory that contains the kubernetes config and the webservice directory that contains the stockticker webapp.

For the webapp, you can run "docker pull geoj/forgerock:release" to pull down the image, then run "docker images to get the id of image. Then run "docker run -p 9000:3000 <image id>" You will then be able to access the webapp at http://localhost:9000.

This was built to run on Intel archs (for reasons I'll discuss later), so if you need to run it on Apple silicon, you can edit the Dockerfile and remove "--platform=linux/amd64". Then docker build . -t <any name you choose>. Then run "docker run -p 9000:3000 <image id>" as above.

The image includes a compiles Go binary, but I've included the Go source code for you to check out.


For the kubernetes config, there were some issues that took me longer to resolve than the actual exercise. What I discoveres is there is an issue with minikube running on Darwin architectures using the docker driver that doesn't allow the host network to connect to the minikube ip. In order to fix that issue, you have to run minikube using another driver, but the only driver supported on Apple silicon is the docker driver.

To get this working, I had to build the image for intel Macs. Then run this on an intel Mac. I installed hyperkit ("brew install hyperkit") deleted my previous install of minikube ("minikube delete") then started minikube with the hyperkit driver ("minikube start --driver=hyperkit")

Then, I had to install the minikube ingress addons:

minikube addons enable ingress
minikube addons enable ingress-dns

Now, get the ip of your minikube cluster with "minikube ip"

You add an entry to the /etc/hosts file to route local traffic to the service ingress, but I added an /etc/resolver/minikube-test file (the resolver directory needed to be created) that looks like this:

domain test
nameserver <minikube ip>
search_order 1
timeout 5

From the k8s directory, run the following:

kubectl apply -f configmap.yaml
kubectl apply -f secrets.yaml
kubectl apply -f pod-service.yaml
kubectl apply -f ingress.yaml

You should be able to access the webapp from a browswer at http://forgerock.test

Once again, this was fun. Hopefully this works for you, I know there are a lot of moving pieces.
