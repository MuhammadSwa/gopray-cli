#!/bin/bash
# Test script to test city preset functionality

echo "Testing City Preset Setup..."

# Remove existing config
rm -f ~/.config/go-pray/conf.yaml

# Test Dubai preset (English interface)
cat << EOF | ./gopray-cli setup
1
2

2
y
EOF

echo ""
echo "Dubai preset setup completed! Testing..."
echo ""
./gopray-cli config
echo ""
./gopray-cli next

echo ""
echo "==============================================="
echo "Testing Cairo preset in Arabic..."
echo ""

# Remove config and test Cairo in Arabic
rm -f ~/.config/go-pray/conf.yaml

cat << EOF | ./gopray-cli setup
2
1


نعم
EOF

echo ""
echo "Cairo Arabic setup completed! Testing..."
echo ""
./gopray-cli config
