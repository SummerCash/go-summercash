cd ..
go build
mv go-summercash dist/go-summercash-$1
rm go-summercash
git add .
git commit -m "Deployed"
git push