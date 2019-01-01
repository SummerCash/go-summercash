# USAGE: ./deploy.sh PLATFORM-NAME (e.g. ./deploy.sh macOS)

cd ..
go build
mv go-summercash dist/go-summercash-$1
git add .
git commit -m "Deployed"
git push