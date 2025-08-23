package chat

import (
	"encoding/json"
	"github.com/pzierahn/chatbot_services/datastore"
	"testing"
)

var threadStr = `{
  "_id": {
    "$binary": {
      "base64": "h51WPAv7QHG151m/Ts0EdA==",
      "subType": "00"
    }
  },
  "collection_id": {
    "$binary": {
      "base64": "WWmHY8D/SMSmnT1q1ip9UA==",
      "subType": "00"
    }
  },
  "messages": [
    {
      "role": "user",
      "content": "Explain how to extract respiration signals with MTTS-CAN"
    },
    {
      "role": "assistant",
      "tool_calls": [
        {
          "id": "e78f566e-2dea-461e-a042-bd513ccb8fcd",
          "name": "attach_document",
          "arguments": "{\"document_id\":\"ecd2821c-fbe9-4026-b4e0-86409a1fa9ab\"}"
        }
      ]
    },
    {
      "role": "user",
      "content": "Rewrite this for a scientific publication. Write without bullet points and lists. Keep the sources and citations!"
    },
    {
      "role": "assistant",
      "content": "The extraction of respiration signals using MTTS-CAN (Multi-Task Temporal Shift Convolutional Attention Network) represents a significant advancement in non-contact physiological measurement. MTTS-CAN employs a sophisticated two-branch architecture incorporating a spatial attention module, as described by Liu et al. \\cite{7b212f00-4554-4ffb-b984-fc280f7f3f0c}. This structure consists of an appearance branch that processes a single averaged frame and a motion branch that utilizes normalized frames derived from adjacent frames in the video sequence.\n\nA key innovation in MTTS-CAN is the implementation of a Temporal Shift Module (TSM) within the motion branch. This module effectively mimics the functionality of 3D convolutions, enabling information exchange among neighboring frames without increasing the parameter count of the network \\cite{5a1eeb6f-fde5-467d-8616-7472ef8d956d}. The architecture also incorporates an attention module that serves as a bridge between the appearance and motion branches. This module generates soft-attention masks, focusing the network's processing on pixels that exhibit stronger physiological signals, thereby enhancing the accuracy of respiration signal extraction \\cite{5a1eeb6f-fde5-467d-8616-7472ef8d956d}.\n\nOne of the distinguishing features of MTTS-CAN is its multi-task learning approach. The network simultaneously estimates both pulse and respiration signals, leveraging shared intermediate representations between these related physiological processes \\cite{5a1eeb6f-fde5-467d-8616-7472ef8d956d}. This approach not only improves computational efficiency but also potentially enhances the accuracy of both measurements through mutual information sharing.\n\nThe respiration signal extraction process culminates in the network outputting a respiration waveform. This waveform undergoes post-processing, which includes the application of a 2nd-order Butterworth filter with cut-off frequencies of 0.08 and 0.5 Hz, specifically tailored for respiration rate estimation \\cite{de6c7968-370e-4d99-91fd-8d442d57acf9}. The efficacy of this approach is evident in the state-of-the-art performance achieved by MTTS-CAN. When evaluated on the AFRL dataset, the network demonstrates a Mean Absolute Error (MAE) of 2.30 breaths per minute for respiration rate estimation \\cite{ae2e2185-8e78-4256-be81-66ec36339e69}.\n\nBeyond its accuracy, MTTS-CAN excels in computational efficiency. The network is capable of processing over 150 frames per second on mobile platforms, with an inference time of only 6ms per frame \\cite{ae2e2185-8e78-4256-be81-66ec36339e69}. This remarkable speed enables real-time applications, making MTTS-CAN particularly suitable for telehealth and remote health monitoring scenarios.\n\nIn conclusion, MTTS-CAN represents a significant leap forward in the extraction of respiration signals from video data. By combining efficient spatial-temporal modeling, attention mechanisms, and multi-task learning, it achieves both high accuracy and real-time performance. These characteristics position MTTS-CAN as a promising tool for various applications in healthcare and beyond, potentially revolutionizing non-contact physiological monitoring."
    },
    {
      "role": "user",
      "content": "Be concise"
    },
    {
      "role": "assistant",
      "content": "MTTS-CAN (Multi-Task Temporal Shift Convolutional Attention Network) extracts respiration signals from video using an innovative architecture. It employs a two-branch structure with a spatial attention module, incorporating a Temporal Shift Module (TSM) for efficient temporal modeling \\cite{7b212f00-4554-4ffb-b984-fc280f7f3f0c}. The network uses attention mechanisms to focus on physiologically relevant pixels and employs multi-task learning to simultaneously estimate pulse and respiration signals \\cite{5a1eeb6f-fde5-467d-8616-7472ef8d956d}. \n\nPost-processing includes a Butterworth filter with specific cut-off frequencies for respiration \\cite{de6c7968-370e-4d99-91fd-8d442d57acf9}. MTTS-CAN achieves state-of-the-art accuracy, with a Mean Absolute Error of 2.30 breaths per minute on the AFRL dataset \\cite{ae2e2185-8e78-4256-be81-66ec36339e69}. It operates at over 150 frames per second on mobile platforms, making it suitable for real-time applications in telehealth and remote monitoring \\cite{ae2e2185-8e78-4256-be81-66ec36339e69}."
    }
  ],
  "timestamp": {
    "$date": "2024-07-13T17:31:38.275Z"
  },
  "user_id": "j7jjxLD9rla2DrZoeUu3Tnft4812"
}`

func Test_messagesToProto(t *testing.T) {
	var thread datastore.Thread
	err := json.Unmarshal([]byte(threadStr), &thread)
	if err != nil {
		t.Error(err)
	}

	out, err := messagesToProto(thread.Messages)
	if err != nil {
		t.Error(err)
	}

	byt, _ := json.MarshalIndent(out, "", "  ")
	t.Log(string(byt))
}
