import { getClient } from "@/lib/client";

type Props = {
  params: Promise<{ id: string }>;
};

export default async function ClientDetails({ params }: Props){
  const { id } = await params;
  const client = await getClient(id);
  console.log(client)


  return (
    <div className="flex flex-row justify-between">
      <div className="flex flex-col">
        <div className="bg-white w-100 p-5 rounded-lg border-2 mt-3">
          <h1 className="text-xl font-bold">{client.name}</h1>
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