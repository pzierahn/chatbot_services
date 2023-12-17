package main

import (
	"context"
	"github.com/pzierahn/brainboost/vertex"
	"log"
)

const text = "1 Introduction\nWhat is simulation? Although it means different things in different settings, there is a clear common denominator. Simulation is a way of comparing what happens in the “real world” to what happens in an “ideal world” where the primitive in question is secure by definition. For example, the definition of semantic security for encryption compares what can be learned by an adversary who receives a real ciphertext to what can be learned by an adversary who receives nothing. The definition states that an encryption scheme is secure if they can both learn approximately the same amount of information. This is very strange. Clearly, the latter adversary who receives nothing can learn nothing about the plaintext since it receives no information. However, this is exactly the point. Since the adversary who receives nothing can learn nothing by triviality (this is an “ideal world” that is secure by definition), this implies that in the real world, where the adversary receives the ciphertext, nothing is learned as well.\nAt first, this seems to be a really complicated way of saying something simple. Why not just define encryption to be secure if nothing is learned? The problem is that it’s not at all clear how to formalize the notion that “nothing is learned”. If we try to say that an adversary who receives a ciphertext cannot output any information about the plaintext, then what happens if the adversary already has information about the plaintext? For example, the adversary may know that it is English text. Of course, this has nothing to do with the security of the scheme since the adversary knew this beforehand and independently of the ciphertext. The simulation-based formulation of security enables us to exactly formalize this. We say that an encryption scheme is secure if the only information derived (or output by the adversary) is that which is based on a priori knowledge. If the adversary receiving no ciphertext is able to output the same information as the adversary receiving the ciphertext, then this is indeed the case.\nIt is unclear at this point why this is called “simulation”; what we have described is a comparison between two worlds. This will be explained throughout the tutorial (first in Section 3). For now, it suffices to say that security proofs for definitions formulated in this way work by constructing a simulator that resides in the alternative world that is secure by definition, and generates a view for the adversary in the real world that is computationally indistinguishable from its real view. In fact, as we will show, there are three distinct but intertwined tasks that a simulator must fulfill:\n1. It must generate a view for the real adversary that is indistinguishable from its real view; 2. It must extract the effective inputs used by the adversary in the execution; and\n3. It must make the view generated be consistent with the output that is based on this input.\nWe will not elaborate on these points here, since it is hard to explain them clearly out of context. However, they will become clear by the end of the tutorial.\nOrganization. In this tutorial, we will demonstrate the simulation paradigm in a number of different settings, together with explanations about what is required from the simulator and proof. We demonstrate the aforementioned three different tasks of the simulator in simulation-based proofs via a gradual progression. Specifically, in Section 3 we provide some more background to the simulation paradigm and how it expresses itself in the context of encryption. Then, in Section 4, we show how to simulate secure computation protocols for the case of semi-honest adversaries (who follow the protocol specification, but try to learn more than allowed by inspecting the protocol\n"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()

	client, err := vertex.New(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}

	embedding, err := client.GenerateEmbeddings(ctx, text)
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("embedding: %v", len(embedding))

	resp, err := client.Generate(ctx, "Sum this up", text)
	if err != nil {
		log.Printf("%v", err)
	}

	log.Printf("resp: %v", resp)
}
