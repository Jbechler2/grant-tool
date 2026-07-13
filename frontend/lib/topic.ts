'use server';

import { Topic } from "@/types";
import { cookies } from "next/headers";

export async function createTopic(label: string): Promise<Topic> {
  const token = (await cookies()).get('token')?.value;

  const url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/topics`
  const res = await fetch(url, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({label})
  });

  if (!res.ok){
    const errorBody = await res.text().catch(() => null);
    throw new Error(`Failed to create topic: ${res.status} ${errorBody ?? ''}`);
  } 
  return res.json();
}

export async function getAllTopics(){ 
  const token = (await cookies()).get('token')?.value;

  const url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/topics`
  const res = await fetch(url, {
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });

  if (!res.ok){
    const errorBody = await res.text().catch(() => null);
    throw new Error(`Failed to get topics: ${res.status} ${errorBody ?? ''}`);
  } 
  return res.json();
}

export async function updateTopic(id: string, label: string){
  const token = (await cookies()).get('token')?.value;

  const url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/topics/${id}`
  const res = await fetch(url, {
    method: 'PUT',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json'
    },  
    body: JSON.stringify({label})
  })

  if (!res.ok) {
    const errorBody = await res.text().catch(() => null)
    throw new Error(`Failed to update topic: ${res.status} ${errorBody ?? ''}`)
  }

  return res.json();
}

export async function deleteTopic(id: string){
  const token = (await cookies()).get('token')?.value;

  const url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/topics/${id}`
  const res = await fetch(url, {
    method: 'DELETE',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  })

  if (!res.ok) {
    const errorBody = await res.text().catch(() => null)
    throw new Error(`Failed to delete topic: ${res.status} ${errorBody ?? ''}`)
  }
}