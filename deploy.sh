v=v0.1.2
git tag $v
git push --tags
go install github.com/ymzuiku/go2ts@$v
echo "done."