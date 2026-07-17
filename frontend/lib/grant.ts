'use server';
import { Grant, Topic } from "@/types";
import { cookies } from "next/headers";

export async function getGrant(id: string): Promise<Grant> {
  const token = (await cookies()).get('token')?.value;

  const url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/grants/${id}`
  const res = await fetch(url, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!res.ok) throw new Error('Failed to fetch grant');
  return res.json();
}

export async function getGrantTopics(id: string): Promise<Topic[]> {
  const token = (await cookies()).get('token')?.value;
  const url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/grants/${id}/topics`

  const res = await fetch(url, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  
  if(!res.ok) throw new Error('Failed to fetch topics');
  return res.json();
}

export async function addTopicToGrant(grantID: string, id: string) {
  const token = (await cookies()).get('token')?.value;

  const url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/grants/${grantID}/topics`
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
}

export async function removeTopicFromGrant(grantID: string, topicID: string) {
  const token = (await cookies()).get('token')?.value;

  const url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/grants/${grantID}/topics/${topicID}`
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