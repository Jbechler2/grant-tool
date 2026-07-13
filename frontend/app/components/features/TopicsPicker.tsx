"use client";

import { Topic } from "@/types";
import { useState } from "react";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import {Command, CommandEmpty, CommandGroup, CommandList, CommandItem, CommandInput} from "@/components/ui/command"
import { Button } from '@/components/ui/button'
import { Check, Plus } from "lucide-react";
import { cn } from "@/lib/utils";
import { createTopic } from "@/lib/topic";

interface TopicSelectionProps {
  parentId: string,
  initialTopics: Topic[],
  allTopics: Topic[]

  removeTopic(parentId: string, topicId: string): Promise<void>;
  addTopic(parentId: string, topicId: string): Promise<void>;
}

export default function TopicsPicker(topicSelectProps: TopicSelectionProps) {
  const [allTopics_copy, setAllTopics_copy] = useState<Topic[]>(topicSelectProps.allTopics)
  const [selectedTopics, setSelectedTopics] = useState<Topic[]>(topicSelectProps.initialTopics);
  const [inputValue, setInputValue] = useState("");
  const [open, setOpen] = useState(false);

  const filteredTopics = topicSelectProps.allTopics.filter(t =>
    t.label.toLowerCase().includes(inputValue.toLowerCase())
  );
  const exactMatch = topicSelectProps.allTopics.some(
    t => t.label.toLowerCase() === inputValue.toLowerCase()
  );

  async function handleSelectTopic(topic: Topic) {
    const alreadySelected = selectedTopics.some(t => t.id == topic.id)
    if(alreadySelected){
      await topicSelectProps.removeTopic(topicSelectProps.parentId, topic.id)
      setSelectedTopics(prev => prev.filter(t => t.id !== topic.id));
    } else {
      await topicSelectProps.addTopic(topicSelectProps.parentId, topic.id)
      setSelectedTopics(prev => [...prev, topic]);
    }

    setInputValue("")
  }

  async function handleCreateTopic() {
    const newTopic = await createTopic(inputValue.trim()); // your API call, returns created topic
    setAllTopics_copy(prev => [...prev, newTopic].sort((a, b) => a.label.localeCompare(b.label)))
    setSelectedTopics(prev => [...prev, newTopic]);
    await topicSelectProps.addTopic(topicSelectProps.parentId, newTopic.id); // link it right away
    setInputValue("");
  }

  return (
    <div>
    <Popover>
      <PopoverTrigger asChild>
        <Button variant="outline">+</Button>
      </PopoverTrigger>
      <PopoverContent align="start" className="w-[400px] p-0">
          <Command shouldFilter={false}>
            <CommandInput placeholder="Search for an existing topic..." onChangeCapture={e => {setInputValue(e.currentTarget.value)}}/>
            <CommandList>
              <CommandEmpty>
                {inputValue.trim() === "" ? "Type to search topics..." : `No matches for "${inputValue.trim()}"`}
              </CommandEmpty>
              <CommandGroup heading="Topics">
                {
                  allTopics_copy.filter(t => t.label.toLowerCase().includes(inputValue.toLowerCase()))
                  .map((topic) => (
                    <CommandItem key={topic.id} onSelect={() => handleSelectTopic(topic)}>
                      <Check className={cn("mr-2 h-4 w-4", selectedTopics.some(t => t.id === topic.id) ? "opacity-100" : "opacity-0")} />
                        {topic.label}
                    </CommandItem>
                  ))
                }
                {inputValue.trim() !== "" && 
                  !allTopics_copy.some(t => t.label.toLowerCase() == inputValue.trim().toLowerCase()) && (
                  <CommandItem onSelect={handleCreateTopic}>
                    <Plus className="mr-2 h-4 w-4" />
                    Create {inputValue.trim()}
                  </CommandItem>
                )}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
    </Popover>
    <div className="flex flex-row gap-2">
      {selectedTopics.map(topic => (
        <h5 key={topic.id} className="bg-green-200 w-fit px-2 py-1 rounded-md text-gray-500">{topic.label}</h5>  
      ))}
    </div>
    </div>
  )
}