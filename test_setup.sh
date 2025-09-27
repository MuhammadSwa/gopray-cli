#!/bin/bash
# Test script to simulate interactive setup

echo "Testing Interactive Setup..."

# Create input for the setup command
cat << EOF | ./gopray-cli setup
1
30.6400
31.2900
Africa/Cairo
7
1
y
EOF

echo ""
echo "Setup completed! Testing the configuration..."
echo ""

# Test the configuration
./gopray-cli config
echo ""
./gopray-cli list
echo ""
./gopray-cli date
