#!/bin/bash
# Generate placeholder audio files using ffmpeg
# These are simple tones to verify the flow works.
# Replace with real recordings later.

AUDIO_DIR="audio"
LANGUAGES=("en" "hi" "mr")
FILES=(
  "welcome:Welcome to Atl Janseva IVR"
  "ward-input:Please enter your pincode and ward"
  "no-match:No matching ward found"
  "ward-menu:Select your ward"
  "nagarsevak-menu:Select your corporator"
  "goodbye:Thank you and goodbye"
  "whatsapp:Please contact us on WhatsApp"
  "main-menu:Main menu options"
  "sos:SOS alert sent"
  "complaint:Complaint registered"
  "corporator-connect:Connecting to corporator"
)

for lang in "${LANGUAGES[@]}"; do
  mkdir -p "$AUDIO_DIR/$lang"
  for entry in "${FILES[@]}"; do
    name="${entry%%:*}"
    text="${entry#*:}"
    ffmpeg -y -f lavfi -i "sine=frequency=440:duration=1" \
      -f lavfi -i "anullsrc=r=24000:cl=mono" \
      -shortest \
      -af "volume=0.3" \
      "$AUDIO_DIR/$lang/$name.mp3" 2>/dev/null
    echo "Created: $AUDIO_DIR/$lang/$name.mp3"
  done
done

echo "Done! Replace these with real recordings when ready."