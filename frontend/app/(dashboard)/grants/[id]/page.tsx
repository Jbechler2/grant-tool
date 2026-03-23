import { getGrant } from "@/lib/grant";

type Props = {
  params: Promise<{ id: string }>;
};

export default async function GrantDetails({ params }: Props){
  const { id } = await params;
  const grant = await getGrant(id);
  console.log(grant)

  // PLACEHOLDER
  const potentialClients = [
    {
      "id": 1,
      "name": "Project Succeed",
    },
    {
      "id": 2,
      "name": "Learning How 2 Live"
    },
    {
      "id": 3,
      "name": "Purses with a Purpose"
    }
  ]
  
  const topics = [
    {
      "id": 1,
      "text": "Community Improvement"
    },
    {
      "id": 2,
      "text": "Basic Needs"
    },
    {
      "id": 3,
      "text": "Workforce Development"
    },
    {
      "id": 4,
      "text": "Arts & Culture"
    },
  ]


  return (
    <div className="flex flex-row w-full justify-between">
      <div>
        <div className="bg-white w-150 p-5 rounded-lg border-2">
          <h1 className="text-xl font-bold">{grant.title}</h1>
          <h3 className="text-sm italic mb-4">{grant.funder_name}</h3>
          <div className="flex flex-row gap-2">
            {topics.map(topic => (
              <h5 key={topic.id} className="bg-green-200 w-fit px-2 py-1 rounded-md text-gray-500">{topic.text}</h5>  
            ))}
          </div>
        </div>
        <div className="bg-white w-100 p-5 rounded-lg border-2 mt-3">
          <h1 className="text-xl font-bold">Grant Amount</h1>
          <h3>{grant.award_amount_max && grant.award_amount_min ? "$" + grant.award_amount_min + " - " + "$" + grant.award_amount_min : (grant.award_amount_max ? "$" + grant.award_amount_max : (grant.award_amount_min ? "$" + grant.award_amount_min : "")) }</h3>
        </div>
        <div className="bg-white w-100 p-5 rounded-lg border-2 mt-3">
          <h1>Deadlines</h1>
          <h3>FETCH DEADLINES</h3>
        </div>
        <div className="bg-white w-100 p-5 rounded-lg border-2 mt-3">
            <h1>Notes</h1>
            <p>{grant.eligibility_notes}</p>
        </div>
      </div>
      <div className="mr-40">
        <div className="bg-white w-100 p-5 rounded-lg border-2 mt-3">
          <h1 className="text-xl font-bold mb-3">Potential Clients</h1>
          {potentialClients.map((client, index) => (
            <div key={client.id} className="">
              {index !== 0 && <hr />}
              <h1>{client.name}</h1>
            </div>  
          ))}
        </div>
        <div className="bg-white w-100 p-5 rounded-lg border-2 mt-3">
          Links
        </div>
      </div>
    </div>
  )
}