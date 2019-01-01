# USAGE: ./deploy.sh PLATFORM-NAME (e.g. ./deploy.sh macOS)

cd ..

go build

mv go-summercash dist/go-summercash-$1

cd dist

if [ "$1" = "macOS" ]
then
    zip go-summercash-$1.zip go-summercash-$1
elif [ "$1" = "win64" ]
then
    zip go-summercash-$1.zip go-summercash-$1.exe
else
    tar -cvzf go-summercash-$1.tar.gz go-summercash-$1
fi

git add .
git commit -m "Deployed"
git push