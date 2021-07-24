for filename in `find . -type f -name 'post*'`; do mv -v "$filename" "${filename//post/drawer}"; done
find ./ -type f \( ! -iname "replace.sh" \) -exec sed -i -e 's/post/drawer/g' {} \;
find ./ -type f \( ! -iname "replace.sh" \) -exec sed -i -e 's/Post/Drawer/g' {} \;
find ./ -type f \( ! -iname "replace.sh" \) -exec sed -i -e 's/MethodDrawer/MethodPost/g' {} \;
find ./ -type f \( ! -iname "replace.sh" \) -exec sed -i -e 's/drawergres/postgres/g' {} \;
find ./ -type f \( ! -iname "replace.sh" \) -exec sed -i -e 's/"blog"/"catwalk"/g' {} \;
find ./ -type f \( ! -iname "replace.sh" \) -exec sed -i -e 's/POSTGRES_DB: blog/POSTGRES_DB: catwalk/g' {} \;
