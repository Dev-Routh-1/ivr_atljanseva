AUDIO FILE RECORDING GUIDE
===========================

Directory structure:
  audio/en/   - English
  audio/hi/   - Hindi
  audio/mr/   - Marathi

Required files and their TTS fallback text:

1. welcome.mp3
   EN: "Welcome to Atl Janseva IVR. Press 1 for English, 2 for Hindi, 3 for Marathi."
   HI: "आटल जनसेवा आईवीआर में आपका स्वागत है। अंग्रेजी के लिए 1 दबाएं, हिंदी के लिए 2, मराठी के लिए 3।"
   MR: "आटल जनसेवा आयव्हीआर मध्ये आपले स्वागत आहे. इंग्रजीसाठी १ दाबा, हिंदीसाठी २, मराठीसाठी ३."

2. ward-input.mp3
   EN: "Please enter your 6 digit pincode followed by hash and your ward name."
   HI: "कृपया अपना 6 अंकों का पिनकोड और वार्ड का नाम हाश के साथ दर्ज करें।"
   MR: "कृपया आपला ६ अंकी पिनकोड आणि वार्डचे नाव हॅशसह प्रविष्ट करा."

3. no-match.mp3
   EN: "We could not find a matching ward. Please try again."
   HI: "हमें कोई मिलता वार्ड नहीं मिला। कृपया पुनः प्रयास करें।"
   MR: "आम्हाला जुळणारा वार्ड सापडला नाही. कृपया पुन्हा प्रयत्न करा."

4. ward-menu.mp3
   EN: "Multiple wards found. Please select your ward."
   HI: "एक से अधिक वार्ड मिले। कृपया अपना वार्ड चुनें।"
   MR: "एकाधिक वार्ड सापडले. कृपया आपला वार्ड निवडा."

5. nagarsevak-menu.mp3
   EN: "Multiple corporators found. Please select one."
   HI: "एक से अधिक नगरसेवक मिले। कृपया एक चुनें।"
   MR: "एकाधिक नगरसेवक सापडले. कृपया एक निवडा."

6. goodbye.mp3
   EN: "Thank you for registering. Goodbye."
   HI: "पंजीकरण के लिए धन्यवाद। नमस्ते।"
   MR: "नोंदणी केल्याबद्दल धन्यवाद. नमस्कार."

7. whatsapp.mp3
   EN: "We could not find your information. Please contact us on WhatsApp for assistance."
   HI: "हमें आपकी जानकारी नहीं मिली। कृपया सहायता के लिए हमें WhatsApp पर संपर्क करें।"
   MR: "आम्हाला तुमची माहिती सापडली नाही. कृपया सहाय्यासाठी आम्हाला WhatsApp वर संपर्क करा."

8. main-menu.mp3
   EN: "Press 1 for SOS, Press 2 to file a complaint, Press 3 to connect to your corporator."
   HI: "SOS के लिए 1 दबाएं, शिकायत दर्ज करने के लिए 2, नगरसेवक से जुड़ने के लिए 3 दबाएं।"
   MR: "SOS साठी १ दाबा, तक्रार नोंदविण्यासाठी २, नगरसेवकाशी संपर्क साधण्यासाठी ३ दाबा."

9. sos.mp3
   EN: "Your SOS alert has been sent. Help is on the way."
   HI: "आपकी SOS सूचना भेज दी गई है। मदद आ रही है।"
   MR: "तुमचा SOS इशारा पाठविला गेला आहे. मदत येत आहे."

10. complaint.mp3
    EN: "Your complaint has been registered. You will receive a response shortly."
    HI: "आपकी शिकायत दर्ज कर ली गई है। जल्द ही आपको जवाब मिलेगा।"
    MR: "तुमची तक्रार नोंदविली गेली आहे. लवकरच तुम्हाला प्रतिसाद मिळेल."

11. corporator-connect.mp3
    EN: "Connecting you to your corporator. Please hold."
    HI: "आपको आपके नगरसेवक से जोड़ा जा रहा है। कृपया प्रतीक्षा करें।"
    MR: "तुम्हाला तुमच्या नगरसेवकाशी जोडले जात आहे. कृपया प्रतीक्षा करा."

Recording specs:
  - Format: MP3
  - Bitrate: 128kbps
  - Sample rate: 44100Hz
  - Channels: Mono
  - Normalize volume to -3dB

Once audio files are ready, update handler/plivo_handler.go to replace
Speak() calls with Play(audioURL) to use your recordings.