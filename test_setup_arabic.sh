#!/bin/bash
# Test script to simulate Arabic interactive setup

echo "Testing Arabic Interactive Setup..."

# Remove existing config
rm -f ~/.config/go-pray/conf.yaml

# Create input for the setup command in Arabic
cat << EOF | ./gopray-cli setup
2
24.4539
54.3773
Asia/Dubai
3
1
نعم
EOF

echo ""
echo "Arabic setup completed! Testing the configuration..."
echo ""

# Test the configuration in Arabic
./gopray-cli config
echo ""
./gopray-cli list
echo ""
./gopray-cli date -d
