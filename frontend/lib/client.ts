'use server'
import { Client, Topic } from "@/types";
import { cookies } from "next/headers";

export async function getClient(id: string): Promise<Client> {
  const token = (await cookies()).get('token')?.value;

  const url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/clients/${id}`
  const res = await fetch(url, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!res.ok) throw new Error('Failed to fetch client');
  return res.json();
}

export async function getClientTopics(id: string): Promise<Topic[]> {
  const token = (await cookies()).get('token')?.value;
  const url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/clients/${id}/topics`

  const res = await fetch(url, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  
  if(!res.ok) throw new Error('Failed to fetch topics');
  return res.json();
}

export async function addTopicToClient(clientID: string, id: string) {
  const token = (await cookies()).get('token')?.value;

  const url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/clients/${clientID}/topics`
  const res = await fetch(url, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({id})
  });

  if (!res.ok){
    const errorBody = await res.text().catch(() => null);
    throw new Error(`Failed to add topic to grant: ${res.status} ${errorBody ?? ''}`);
  } 
  return res.json();
}

export async function removeTopicFromClient(clientID: string, topicID: string) {
  const token = (await cookies()).get('token')?.value;

  const url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/clients/${clientID}/topics/${topicID}`
  const res = await fetch(url, {
    method: 'DELETE',
    headers: {
      Authorization: `Bearer ${token}`,
    }
  });

  if (!res.ok){
    const errorBody = await res.text().catch(() => null);
    throw new Error(`Failed to remove topic from grant: ${res.status} ${errorBody ?? ''}`);
  }
}