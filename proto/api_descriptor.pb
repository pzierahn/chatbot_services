
�
google/protobuf/empty.protogoogle.protobuf"
EmptyB}
com.google.protobufB
EmptyProtoPZ.google.golang.org/protobuf/types/known/emptypb��GPB�Google.Protobuf.WellKnownTypesJ�
 2
�
 2� Protocol Buffers - Google's data interchange format
 Copyright 2008 Google Inc.  All rights reserved.
 https://developers.google.com/protocol-buffers/

 Redistribution and use in source and binary forms, with or without
 modification, are permitted provided that the following conditions are
 met:

     * Redistributions of source code must retain the above copyright
 notice, this list of conditions and the following disclaimer.
     * Redistributions in binary form must reproduce the above
 copyright notice, this list of conditions and the following disclaimer
 in the documentation and/or other materials provided with the
 distribution.
     * Neither the name of Google Inc. nor the names of its
 contributors may be used to endorse or promote products derived from
 this software without specific prior written permission.

 THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


  

" E
	
" E

# ,
	
# ,

$ +
	
$ +

% "
	

% "

& !
	
$& !

' ;
	
%' ;

( 
	
( 
�
 2 � A generic empty message that you can re-use to avoid defining duplicated
 empty messages in your APIs. A typical example is to use it as the request
 or the response type of an API method. For instance:

     service Foo {
       rpc Bar(google.protobuf.Empty) returns (google.protobuf.Empty);
     }




 2bproto3
�1
google/protobuf/timestamp.protogoogle.protobuf";
	Timestamp
seconds (Rseconds
nanos (RnanosB�
com.google.protobufBTimestampProtoPZ2google.golang.org/protobuf/types/known/timestamppb��GPB�Google.Protobuf.WellKnownTypesJ�/
 �
�
 2� Protocol Buffers - Google's data interchange format
 Copyright 2008 Google Inc.  All rights reserved.
 https://developers.google.com/protocol-buffers/

 Redistribution and use in source and binary forms, with or without
 modification, are permitted provided that the following conditions are
 met:

     * Redistributions of source code must retain the above copyright
 notice, this list of conditions and the following disclaimer.
     * Redistributions in binary form must reproduce the above
 copyright notice, this list of conditions and the following disclaimer
 in the documentation and/or other materials provided with the
 distribution.
     * Neither the name of Google Inc. nor the names of its
 contributors may be used to endorse or promote products derived from
 this software without specific prior written permission.

 THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


  

" 
	
" 

# I
	
# I

$ ,
	
$ ,

% /
	
% /

& "
	

& "

' !
	
$' !

( ;
	
%( ;
�
 � �� A Timestamp represents a point in time independent of any time zone or local
 calendar, encoded as a count of seconds and fractions of seconds at
 nanosecond resolution. The count is relative to an epoch at UTC midnight on
 January 1, 1970, in the proleptic Gregorian calendar which extends the
 Gregorian calendar backwards to year one.

 All minutes are 60 seconds long. Leap seconds are "smeared" so that no leap
 second table is needed for interpretation, using a [24-hour linear
 smear](https://developers.google.com/time/smear).

 The range is from 0001-01-01T00:00:00Z to 9999-12-31T23:59:59.999999999Z. By
 restricting to that range, we ensure that we can convert to and from [RFC
 3339](https://www.ietf.org/rfc/rfc3339.txt) date strings.

 # Examples

 Example 1: Compute Timestamp from POSIX `time()`.

     Timestamp timestamp;
     timestamp.set_seconds(time(NULL));
     timestamp.set_nanos(0);

 Example 2: Compute Timestamp from POSIX `gettimeofday()`.

     struct timeval tv;
     gettimeofday(&tv, NULL);

     Timestamp timestamp;
     timestamp.set_seconds(tv.tv_sec);
     timestamp.set_nanos(tv.tv_usec * 1000);

 Example 3: Compute Timestamp from Win32 `GetSystemTimeAsFileTime()`.

     FILETIME ft;
     GetSystemTimeAsFileTime(&ft);
     UINT64 ticks = (((UINT64)ft.dwHighDateTime) << 32) | ft.dwLowDateTime;

     // A Windows tick is 100 nanoseconds. Windows epoch 1601-01-01T00:00:00Z
     // is 11644473600 seconds before Unix epoch 1970-01-01T00:00:00Z.
     Timestamp timestamp;
     timestamp.set_seconds((INT64) ((ticks / 10000000) - 11644473600LL));
     timestamp.set_nanos((INT32) ((ticks % 10000000) * 100));

 Example 4: Compute Timestamp from Java `System.currentTimeMillis()`.

     long millis = System.currentTimeMillis();

     Timestamp timestamp = Timestamp.newBuilder().setSeconds(millis / 1000)
         .setNanos((int) ((millis % 1000) * 1000000)).build();

 Example 5: Compute Timestamp from Java `Instant.now()`.

     Instant now = Instant.now();

     Timestamp timestamp =
         Timestamp.newBuilder().setSeconds(now.getEpochSecond())
             .setNanos(now.getNano()).build();

 Example 6: Compute Timestamp from current time in Python.

     timestamp = Timestamp()
     timestamp.GetCurrentTime()

 # JSON Mapping

 In JSON format, the Timestamp type is encoded as a string in the
 [RFC 3339](https://www.ietf.org/rfc/rfc3339.txt) format. That is, the
 format is "{year}-{month}-{day}T{hour}:{min}:{sec}[.{frac_sec}]Z"
 where {year} is always expressed using four digits while {month}, {day},
 {hour}, {min}, and {sec} are zero-padded to two digits each. The fractional
 seconds, which can go up to 9 digits (i.e. up to 1 nanosecond resolution),
 are optional. The "Z" suffix indicates the timezone ("UTC"); the timezone
 is required. A proto3 JSON serializer should always use UTC (as indicated by
 "Z") when printing the Timestamp type and a proto3 JSON parser should be
 able to accept both UTC and other timezones (as indicated by an offset).

 For example, "2017-01-15T01:30:15.01Z" encodes 15.01 seconds past
 01:30 UTC on January 15, 2017.

 In JavaScript, one can convert a Date object to this format using the
 standard
 [toISOString()](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date/toISOString)
 method. In Python, a standard `datetime.datetime` object can be converted
 to this format using
 [`strftime`](https://docs.python.org/2/library/time.html#time.strftime) with
 the time format spec '%Y-%m-%dT%H:%M:%S.%fZ'. Likewise, in Java, one can use
 the Joda Time's [`ISODateTimeFormat.dateTime()`](
 http://joda-time.sourceforge.net/apidocs/org/joda/time/format/ISODateTimeFormat.html#dateTime()
 ) to obtain a formatter capable of generating timestamps in this format.



 �
�
  �� Represents seconds of UTC time since Unix epoch
 1970-01-01T00:00:00Z. Must be from 0001-01-01T00:00:00Z to
 9999-12-31T23:59:59Z inclusive.


  �

  �

  �
�
 �� Non-negative fractions of a second at nanosecond resolution. Negative
 second values with fractions must still have non-negative nanos values
 that count forward in time. Must be from 0 to 999,999,999
 inclusive.


 �

 �

 �bproto3
�
account_service.protochatbot.account.v1google/protobuf/empty.protogoogle/protobuf/timestamp.proto"�
BalanceSheet7
payments (2.chatbot.account.v1.PaymentRpayments4
costs (2.chatbot.account.v1.ModelCostsRcosts
balance (Rbalance"�

ModelCosts
model (	Rmodel
input (Rinput
output (Routput
costs (Rcosts
requests (Rrequests"?
Costs6
models (2.chatbot.account.v1.ModelCostsRmodels"a
Payment
id (	Rid.
date (2.google.protobuf.TimestampRdate
amount (Ramount"=
Payments1
items (2.chatbot.account.v1.PaymentRitems2�
AccountService=
GetCosts.google.protobuf.Empty.chatbot.account.v1.CostsC
GetPayments.google.protobuf.Empty.chatbot.account.v1.PaymentsK
GetBalanceSheet.google.protobuf.Empty .chatbot.account.v1.BalanceSheetB	Z./protoJ�
  *

  

 
	
 

 
	
  %
	
 )


 	 


 	

  
6

  


  
$

  
/4

 <

 

 '

 2:

 D

 

 +

 6B


  


 

   

  


  

  

  

  

 


 

 

 

 

 

 

 


 




 

 

 	

 





	







	







	







	




  




 !

 


 

 

  


" &


"

 #

 #

 #	

 #

$%

$

$ 

$#$

%

%

%	

%


( *


(

 )

 )


 )

 )

 )bproto3
�
collection_service.protochatbot.collections.v2google/protobuf/empty.proto"
CollectionID
id (	Rid"V

Collection
id (	Rid
name (	Rname$
documentCount (RdocumentCount"G
Collections8
items (2".chatbot.collections.v2.CollectionRitems2�
CollectionServiceO
Get$.chatbot.collections.v2.CollectionID".chatbot.collections.v2.CollectionC
List.google.protobuf.Empty#.chatbot.collections.v2.CollectionsP
Create".chatbot.collections.v2.Collection".chatbot.collections.v2.CollectionP
Update".chatbot.collections.v2.Collection".chatbot.collections.v2.CollectionD
Delete".chatbot.collections.v2.Collection.google.protobuf.EmptyB	Z./protoJ�
  

  

 
	
 

 
	
  %


  


 

  	-

  		

  	


  	!+

 
8

 



 
 

 
+6

 .

 

 

 ",

 .

 

 

 ",

 9

 

 

 "7


  


 

  

  

  	

  


 




 

 

 	

 





	







	




 




  

 


 

 

 bproto3
�+
chat_service.protochatbot.chat.v4google/protobuf/timestamp.protogoogle/protobuf/empty.protocollection_service.proto"�
CompletionRequest
document_id (	R
documentId
prompt (	RpromptB
model_options (2.chatbot.chat.v4.ModelOptionsRmodelOptions"4
CompletionResponse

completion (	R
completion"�
BatchRequest!
document_ids (	RdocumentIds
prompts (	RpromptsB
model_options (2.chatbot.chat.v4.ModelOptionsRmodelOptions"�
BatchResponse!
document_ids (	RdocumentIds
prompts (	Rprompts!
prompt_title (	RpromptTitle?
items (2).chatbot.chat.v4.BatchResponse.CompletionRitems�

Completion
document_id (R
documentId%
document_title (	RdocumentTitle
prompt (Rprompt

completion (	R
completion"�
Prompt
threadID (	RthreadID
prompt (	RpromptB
model_options (2.chatbot.chat.v4.ModelOptionsRmodelOptions"�
ThreadPrompt
prompt (	Rprompt#
collection_id (	RcollectionIdB
model_options (2.chatbot.chat.v4.ModelOptionsRmodelOptions
	threshold (R	threshold
limit (Rlimit!
document_ids (	RdocumentIds"z
ModelOptions
model (	Rmodel 
temperature (Rtemperature

max_tokens (R	maxTokens
top_p (RtopP"�
Message
id (	Rid
prompt (	Rprompt

completion (	R
completion8
	timestamp (2.google.protobuf.TimestampR	timestamp"�
Thread
id (	Rid4
messages (2.chatbot.chat.v4.MessageRmessages"
referenceIDs (	RreferenceIDsW
reference_scores (2,.chatbot.chat.v4.Thread.ReferenceScoresEntryRreferenceScoresB
ReferenceScoresEntry
key (	Rkey
value (Rvalue:8"
ThreadID
id (	Rid"8
	MessageID
id (	Rid
	thread_id (	RthreadId"
	ThreadIDs
ids (	Rids2�
ChatServiceE
StartThread.chatbot.chat.v4.ThreadPrompt.chatbot.chat.v4.Thread@
PostMessage.chatbot.chat.v4.Prompt.chatbot.chat.v4.Message?
	GetThread.chatbot.chat.v4.ThreadID.chatbot.chat.v4.ThreadO
ListThreadIDs".chatbot.collections.v2.Collection.chatbot.chat.v4.ThreadIDsA
DeleteThread.chatbot.chat.v4.ThreadID.google.protobuf.EmptyM
DeleteMessageFromThread.chatbot.chat.v4.MessageID.google.protobuf.EmptyU

Completion".chatbot.chat.v4.CompletionRequest#.chatbot.chat.v4.CompletionResponseB	Z./protoJ�
  e

  

 
	
 

 
	
  )
	
 %
	
	 "


  


 

  1

  

  

  )/

 ,

 

 

 #*

 +

 

 

 #)

 C

 

 -

 8A

 =

 

 

 &;

 I

 

 '

 2G

 A

 

 "

 -?


  


 

  

  

  	

  

 

 

 	

 

 !

 

 

  


 




 

 

 	

 


 #




  #

  


  

  

  !"

!

!


!

!

!

"!

"

"

" 


% 2


%

 &#

 &


 &

 &

 &!"

'

'


'

'

'

(#

(


(

(

(!"

 */

 *


  +

  +


  +

  +

 ,

 ,


 ,

 ,

 -

 -


 -

 -

 .

 .


 .

 .

1 

1


1

1

1


3 7


3

 4

 4

 4	

 4

5

5

5	

5

6!

6

6

6 


9 C


9

 :

 :

 :	

 :

;

;

;	

;

<!

<

<

< 

? Search options


?

?

?

@

@

@	

@

B#

B


B

B

B!"


E J


E

 F

 F

 F	

 F

G

G

G

G

H

H

H	

H

I

I

I

I


L Q


L

 M

 M

 M	

 M

N

N

N	

N

O

O

O	

O

P*

P

P%

P()


S X


S

 T

 T

 T	

 T

U 

U


U

U

U

V#

V


V

V

V!"

W*

W

W%

W()


	Z \


	Z

	 [

	 [

	 [	

	 [



^ a



^


 _


 _


 _	


 _


`


`


`	


`


c e


c

 d

 d


 d

 d

 dbproto3
�
crashlytics.protocrashlytics.v1google/protobuf/empty.proto"g
Error
	exception (	R	exception
stack_trace (	R
stackTrace
app_version (	R
appVersion2R
CrashlyticsService<
RecordError.crashlytics.v1.Error.google.protobuf.EmptyB	Z./protoJ�
  

  

 
	
 

 
	
  %


  



 

  	9

  	

  	

  	"7


  


 

  

  

  	

  

 

 

 	

 

 

 

 	

 bproto3
�/
document_service.protochatbot.documents.v2google/protobuf/empty.protogoogle/protobuf/timestamp.protocollection_service.proto"�
DocumentNamesD
items (2..chatbot.documents.v2.DocumentNames.ItemsEntryRitems8

ItemsEntry
key (	Rkey
value (	Rvalue:8"s
RenameDocument
id (	Rid
	file_name (	H RfileName%
webpage_title (	H RwebpageTitleB
	rename_to"

DocumentID
id (	Rid"�
DocumentListC
items (2-.chatbot.documents.v2.DocumentList.ItemsEntryRitems`

ItemsEntry
key (	Rkey<
value (2&.chatbot.documents.v2.DocumentMetadataRvalue:8"$
ReferenceIDs
items (	Ritems"A
Chunk
id (	Rid
text (	Rtext
index (Rindex"B

References4
items (2.chatbot.documents.v2.DocumentRitems"|
SearchQuery
query (	Rquery#
collection_id (	RcollectionId
	threshold (R	threshold
limit (Rlimit"�
SearchResults4
items (2.chatbot.documents.v2.DocumentRitemsG
scores (2/.chatbot.documents.v2.SearchResults.ScoresEntryRscores9
ScoresEntry
key (	Rkey
value (Rvalue:8"C
IndexProgress
status (	Rstatus
progress (Rprogress"K
DocumentFilter
query (	Rquery#
collection_id (	RcollectionId"
DocumentMetadata0
file (2.chatbot.documents.v2.FileH Rfile1
web (2.chatbot.documents.v2.WebpageH RwebB
data"6
File
path (	Rpath
filename (	Rfilename"1
Webpage
url (	Rurl
title (	Rtitle"�
Document
id (	Rid#
collection_id (	RcollectionId9

created_at (2.google.protobuf.TimestampR	createdAtB
metadata (2&.chatbot.documents.v2.DocumentMetadataRmetadata3
chunks (2.chatbot.documents.v2.ChunkRchunks"�
DocumentHeader
id (	Rid#
collection_id (	RcollectionId9

created_at (2.google.protobuf.TimestampR	createdAtB
metadata (2&.chatbot.documents.v2.DocumentMetadataRmetadata"�
IndexJob
id (	Rid#
collection_id (	RcollectionIdB
document (2&.chatbot.documents.v2.DocumentMetadataRdocument2�
DocumentServiceP
List$.chatbot.documents.v2.DocumentFilter".chatbot.documents.v2.DocumentListG
Get .chatbot.documents.v2.DocumentID.chatbot.documents.v2.DocumentS
	GetHeader .chatbot.documents.v2.DocumentID$.chatbot.documents.v2.DocumentHeaderF
Rename$.chatbot.documents.v2.RenameDocument.google.protobuf.EmptyB
Delete .chatbot.documents.v2.DocumentID.google.protobuf.EmptyN
Index.chatbot.documents.v2.IndexJob#.chatbot.documents.v2.IndexProgress0P
Search!.chatbot.documents.v2.SearchQuery#.chatbot.documents.v2.SearchResultsU
GetReferences".chatbot.documents.v2.ReferenceIDs .chatbot.documents.v2.References]
MapDocumentNames$.chatbot.collections.v2.CollectionID#.chatbot.documents.v2.DocumentNamesB	Z./protoJ�
  t

  

 
	
 

 
	
  %
	
 )
	
	 "


  


 

  2

  


  

  $0

 )

 	

 


 '

 5

 

 

 %3

 =

 

 

 &;

 9

 

 

 "7

 5

 

 

 %

 &3

 2

 

 

 #0

 7

 

  

 +5

 T

 

 :

 ER


  


 

   

  

  

  


 !




 

 

 	

 

  

 




















# %


#

 $

 $

 $	

 $


' *


'

 )* Id to filename


 )

 ) %

 )()


, .


,

 -

 -


 -

 -

 -


0 4


0

 1

 1

 1	

 1

2

2

2	

2

3

3

3	

3


6 8


6

 7

 7


 7

 7

 7


: ?


:

 ;

 ;

 ;	

 ;

<

<

<	

<

=

=

=

=

>

>

>	

>


A D


A

 B

 B


 B

 B

 B

C 

C

C

C


	F I


	F

	 G

	 G

	 G	

	 G

	H

	H

	H

	H



K N



K


 L


 L


 L	


 L


M


M


M	


M


P U


P

 QT

 Q

 R

 R

 R	

 R

S

S

S

S


W Z


W

 X

 X

 X	

 X

Y

Y

Y	

Y


\ _


\

 ]

 ]

 ]	

 ]

^

^

^	

^


a g


a

 b

 b

 b	

 b

c

c

c	

c

d+

d

d&

d)*

e 

e

e

e

f

f


f

f

f


i n


i

 j

 j

 j	

 j

k

k

k	

k

l+

l

l&

l)*

m 

m

m

m


p t


p

 q

 q

 q	

 q

r

r

r	

r

s 

s

s

sbproto3
�
notion.protochatbot.notion.v1google/protobuf/empty.protochat_service.proto" 
NotionApiKey
key (	Rkey"�
NotionPrompt

databaseID (	R
databaseID"
collectionID (	RcollectionID
prompt (	RpromptA
modelOptions (2.chatbot.chat.v4.ModelOptionsRmodelOptions"-
ExecutionResult
document (	Rdocument"
DatabasesID
id (	Rid"p
	Databases7
items (2!.chatbot.notion.v1.Databases.ItemRitems*
Item
id (	Rid
name (	Rname2�
NotionD
	SetApiKey.chatbot.notion.v1.NotionApiKey.google.protobuf.Empty>
RemoveApiKey.google.protobuf.Empty.google.protobuf.EmptyD
	GetApiKey.google.protobuf.Empty.chatbot.notion.v1.NotionApiKeyE
ListDatabases.google.protobuf.Empty.chatbot.notion.v1.DatabasesV
ExecutePrompt.chatbot.notion.v1.NotionPrompt".chatbot.notion.v1.ExecutionResult0B	Z./protoJ�
  ,

  

 
	
 

 
	
  %
	
 


 
 


 


  >

  

  

  '<

 J

 

 (

 3H

 >

 

 %

 0<

 ?

 

 )

 4=

 C

 

  

 +1

 2A


  


 

  

  

  	

  


 




 

 

 	

 





	







	



(



#

&'


 




 

 

 	

 


! #


!

 "

 "

 "	

 "


% ,


%

 &)

 &


  '

  '


  '

  '

 (

 (


 (

 (

 +

 +


 +

 +

 +bproto3