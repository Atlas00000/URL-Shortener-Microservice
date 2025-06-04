#!/bin/bash

# Load environment variables if .env file exists
if [ -f .env ]; then
    source .env
fi

# Check if required environment variables are set
if [ -z "$MAXMIND_ACCOUNT_ID" ] || [ -z "$MAXMIND_LICENSE_KEY" ]; then
    echo "Error: MAXMIND_ACCOUNT_ID and MAXMIND_LICENSE_KEY must be set in .env file"
    exit 1
fi

# Create directories if they don't exist
mkdir -p data/geoip

# Download GeoLite2 City database
curl -s "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=${MAXMIND_LICENSE_KEY}&suffix=tar.gz" \
    -o data/geoip/GeoLite2-City.tar.gz

# Extract the database
tar -xzf data/geoip/GeoLite2-City.tar.gz -C data/geoip

# Move the .mmdb file to the correct location
mv data/geoip/GeoLite2-City_*/GeoLite2-City.mmdb data/geoip/

# Clean up
rm -rf data/geoip/GeoLite2-City_* data/geoip/GeoLite2-City.tar.gz

echo "GeoLite2 City database downloaded and extracted successfully" 