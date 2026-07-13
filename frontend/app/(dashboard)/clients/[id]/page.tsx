import TopicsPicker from "@/app/components/features/TopicsPicker";
import { addTopicToClient, getClient, getClientTopics, removeTopicFromClient } from "@/lib/client";
import { getAllTopics } from "@/lib/topic";

type Props = {
  params: Promise<{ id: string }>;
};

export default async function ClientDetails({ params }: Props){
  const { id } = await params;
  const client = await getClient(id);
  const grant_topics = await getClientTopics(id)
  const all_topics = await getAllTopics()
  

  return (
    <div className="flex flex-row justify-between">
      <div className="flex flex-col">
        <div className="bg-white w-100 p-5 rounded-lg border-2 mt-3">
          <h1 className="text-xl font-bold">{client.name}</h1>
          <TopicsPicker parentId={ id } initialTopics={ grant_topics } allTopics={ all_topics } addTopic={addTopicToClient} removeTopic={removeTopicFromClient}></TopicsPicker>
        </div>
        <div className="bg-white w-100 p-5 rounded-lg border-2 mt-3">
          <h1 className="text-xl font-bold">Contact Info</h1>
          <div className="flex flex-row gap-4">
            <label className="font-bold">Primary Contact: </label>
            <h3>{client.contact_name}</h3>
          </div>
          <div className="flex flex-row gap-4">
            <label className="font-bold">Contact Phone #: </label>
            <h3>{client.contact_phone}</h3>
          </div>
          <div className="flex flex-row gap-4">
            <label className="font-bold">Contact Email: </label>
            <h3>{client.contact_email}</h3>
          </div>
        </div>
      </div>
      <div>
        <div className="bg-white w-100 p-5 rounded-lg border-2 mt-3">
          <h1 className="text-xl font-bold">Associated Grants</h1>
        </div>
      </div>
    </div>
  )
}