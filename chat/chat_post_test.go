package chat

import (
	"encoding/json"
	"github.com/pzierahn/chatbot_services/datastore"
	"log"
	"testing"
)

var threadStr = `{
  "_id": {
    "$binary": {
      "base64": "cAhj1aLUQreXc3bauelAAg==",
      "subType": "00"
    }
  },
  "collection_id": {
    "$binary": {
      "base64": "DwpB51BpSW6vx8RR5gdj9Q==",
      "subType": "00"
    }
  },
  "messages": [
    {
      "role": "user",
      "content": "Explain pascal"
    },
    {
      "role": "assistant",
      "tool_calls": [
        {
          "tool_call_id": "toolu_bdrk_01Ju4ddfGr41dDEeC6AG1tBC",
          "function": {
            "name": "get_sources",
            "arguments": "{\"prompt\":\"pascal programming language\"}"
          }
        }
      ]
    },
    {
      "role": "user",
      "tool_responses": [
        {
          "tool_call_id": "toolu_bdrk_01Ju4ddfGr41dDEeC6AG1tBC",
          "content": "{\"sources\":[{\"id\":\"3adc335e-03ce-4fe2-ad1b-25a47ab04a0c\",\"text\":\"3\\nPARADIGM OVERVIEW\\n\\nThe three paradigms included in this overview chapter are structured programming, object-orient\\nprogramming, and functional programming.\\n\\nSTRUCTURED PROGRAMMING\\nThe first paradigm to be adopted (but not the first to be invented) was structured programming, which\\nwas discovered by Edsger Wybe Dijkstra in 1968. Dijkstra showed that the use of unrestrained jumps\\n(goto statements) is harmful to program structure. As we’ll see in the chapters that follow, he\\nreplaced those jumps with the more familiar if/then/else and do/while/until constructs.\\nWe can summarize the structured programming paradigm as follows:\\nStructured programming imposes discipline on direct transfer of control.\\n\\nOBJECT-ORIENTED PROGRAMMING\\nThe second paradigm to be adopted was actually discovered two years earlier, in 1966, by Ole Johan\\nDahl and Kristen Nygaard. These two programmers noticed that the function call stack frame in the\\nALGOL language could be moved to a heap, thereby allowing local variables declared by a function to\",\"document_id\":\"7af7b28a-08f0-4dc3-a0a6-a99f7205fd13\",\"position\":36,\"score\":0.40957278},{\"id\":\"194bd103-5de0-49d7-810d-31d8c7d029e4\",\"text\":\"Click here to view code image\\n\\nSTART,\\n\\nBUFR,\\nGETSTR,\\nNXTCH,\\n\\nK177,\\nMCR,\\n\\n*200\\nTLS\\nCLA\\nTAD BUFR\\nJMS GETSTR\\nCLA\\nTAD BUFR\\nJMS PUTSTR\\nJMP START\\n3000\\n0\\nDCA PTR\\nKSF\\nJMP -1\\nKRB\\nDCA I PTR\\nTAD I PTR\\nAND K177\\nISZ PTR\\nTAD MCR\\nSZA\\nJMP NXTCH\\n177\\n-15\\n\\nNote the *200 command at the start of this program. It tells the compiler to generate code that will be\\nloaded at address 2008.\\nThis kind of programming is a foreign concept for most programmers today. They rarely have to think\\nabout where a program is loaded in the memory of the computer. But in the early days, this was one of\\nthe first decisions a programmer needed to make. In those days, programs were not relocatable.\\nHow did you access a library function in those olden days? The preceding code illustrates the\\napproach used. Programmers included the source code of the library functions with their application\\ncode, and compiled them all as a single program.1 Libraries were kept in source, not in binary.\\nThe problem with this approach was that, during this era, devices were slow and memory was\\nexpensive and, therefore, limited. Compilers needed to make several passes over the source code, but\\nmemory was too limited to keep all the source code resident. Consequently, the compiler had to read\\nin the source code several times using the slow devices.\\nThis took a long time—and the larger your function library, the longer the compiler took. Compiling a\\nlarge program could take hours.\\nTo shorten the compile times, programmers separated the source code of the function library from the\\napplications. They compiled the function library separately and loaded the binary at a known address\\n—say, 20008. They created a symbol table for the function library and compiled that with their\\napplication code. When they wanted to run an application, they would load the binary function\",\"document_id\":\"7af7b28a-08f0-4dc3-a0a6-a99f7205fd13\",\"position\":87,\"score\":0.3836394},{\"id\":\"ac2b25ad-1341-4a80-b8e0-08bd6b218131\",\"text\":\"convenient language with which to program the machine. The 8085 assembler was simply not fun to\\nwrite code in.\\nOn top of that, the assembler we were using was written by our own programmers. It ran on our\\nM365 computers, using the cartridge tape operating system described in the “Laser Trim” section.\\nAs fate would have it, our lead hardware engineer convinced our CEO that we needed a real\\ncomputer. He didn’t actually know what he would do with it, but he had a lot of political clout. So we\\npurchased a PDP-11/60.\\nI, a lowly software developer at the time, was ecstatic. I knew precisely what I wanted to do with\\nthat computer. I was determined that this was going to be my machine.\\nWhen the manuals arrived, many months before the delivery of the machine, I took them home and\\ndevoured them. By the time the computer was delivered, I knew how to operate both the hardware\\nand the software at an intimate level—at least, as intimate as home study can make it.\\nI helped to write the purchase order. In particular, I specified the disk storage that the new computer\\nwould have. I decided we should buy two disk drives that could take removable disk packs that held\\n25 megabytes each.9\\nFifty megabytes! The number seemed infinite! I remember walking through the halls of the office, late\\nat night, cackling like the Wicked Witch of the West: “Fifty megabytes! Hahahahahahahahahah!”\\nI had the facilities manager build a little room that would house six VT100 terminals. I decorated it\\nwith pictures from space. Our software developers would use this room to write and compile code.\\nWhen the machine arrived, I spent several days setting it up, wiring all the terminals, and getting\\neverything to work. It was a joy—a labor of love.\\nWe purchased standard assemblers for the 8085 from Boston Systems Office, and we translated the 4Tel Micro code into that syntax. We built a cross-compilation system that allowed us to download\\ncompiled binaries from the PDP-11 to our 8085 development environments, and ROM burners. And\\n—Bob’s your Uncle—it all worked like a champ.\\n\\nC\\nBut that left us with the problem of still using 8085 assembler. That was not a situation that I was\\nhappy with. I had heard that there was this “new” language that was heavily used at Bell Labs. They\\ncalled it “C.” So I purchased a copy of The C Programming Language by Kernighan and Ritchie.\\nLike the PDP-11 manuals a few months before, I inhaled this book.\\nI was astounded by the simple elegance of this language. It sacrificed none of the power of assembly\\nlanguage, and provided access to that power with a much more convenient syntax. I was sold.\\nI purchased a C compiler from Whitesmiths, and got it running on the PDP-11. The output of the\",\"document_id\":\"7af7b28a-08f0-4dc3-a0a6-a99f7205fd13\",\"position\":256,\"score\":0.36392978}]}"
        }
      ]
    },
    {
      "role": "assistant",
      "content": "Pascal is a high-level programming language that was developed in the late 1960s by Niklaus Wirth. Some key points about Pascal:\n\n1. Structured Programming: Pascal was designed to promote structured programming, which emphasizes the use of control structures like if-then-else and while loops instead of unrestrained jumps (goto statements). This helps improve the readability and maintainability of code.\n\n2. Emphasis on Data Types: Pascal has a strong emphasis on data types, requiring variables to be declared with a specific type. This helps catch type-related errors at compile-time.\n\n3. Modular Design: Pascal encourages the use of modular program design through the use of procedures and functions. This helps organize large programs into reusable components.\n\n4. Educational Language: Pascal was originally designed as an educational language to teach programming concepts. It was widely used in universities and schools to introduce students to programming.\n\n5. Influence on Other Languages: Pascal had a significant influence on the design of many other programming languages, including C, Ada, and Delphi. Many of its structured programming concepts became standard in modern languages.\n\n6. Applications: While Pascal is not as widely used today as languages like C, C++, or Java, it is still used in certain domains, such as systems programming, embedded systems, and education.\n\nIn summary, Pascal was an influential programming language that helped establish structured programming paradigms and data typing as key principles of modern programming languages. Its focus on readability and modular design made it a popular choice for teaching programming concepts."
    }
  ],
  "timestamp": {
    "$date": "2024-07-09T10:28:53.334Z"
  },
  "user_id": "j7jjxLD9rla2DrZoeUu3Tnft4812"
}`

func Test_getSources(t *testing.T) {
	var thread datastore.Thread
	_ = json.Unmarshal([]byte(threadStr), &thread)

	sources := getSources(thread.Messages)
	log.Printf("######### Sources: %v", sources)
}
