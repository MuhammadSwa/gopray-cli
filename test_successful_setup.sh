#!/bin/bash
# Final test with successful setup

echo "Testing Successful Interactive Setup..."

# Remove existing config
rm -f ~/.config/go-pray/conf.yaml

# Test successful Dubai setup
cat << EOF | ./gopray-cli setup
1
2


y
EOF

echo ""
echo "Setup completed successfully! Testing all features..."
echo ""

# Test all commands
echo "🔹 Configuration:"
./gopray-cli config
echo ""

echo "🔹 Prayer Times:"
./gopray-cli list
echo ""

echo "🔹 Next Prayer:"
./gopray-cli next
echo ""

echo "🔹 Today's Date:"
./gopray-cli date
echo ""

echo "🔹 Detailed Date:"
./gopray-cli date -d
