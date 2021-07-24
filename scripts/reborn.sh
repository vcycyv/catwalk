sed -i 's/\[replacement\]/drawer/g' replace.sh
sed -i 's/\[Replacement\]/Drawer/g' replace.sh
sed -i 's/\[project\]/catwalk/g' replace.sh
chmod +x replace.sh
mv replace.sh ..
cd ..
./replace.sh
mv replace.sh scripts/
cd scripts