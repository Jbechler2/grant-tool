"use client"
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { getAllTopics, deleteTopic, updateTopic } from "@/lib/topic";
import { Topic } from "@/types";
import { Trash } from "lucide-react";
import { useEffect, useState } from "react";

interface EditTopicsProps {
  open: boolean,
  onOpenChange: (open: boolean) => void
}

export default function EditTopics({open, onOpenChange}: EditTopicsProps) {
  const [topics, setTopics] = useState<Topic[]>([]);
  const [loading, setLoading] = useState<boolean>(false)

  const [editingId, setEditingId] = useState<string | null>(null);
  const [editValue, setEditValue] = useState("");
  const [deleteTarget, setDeleteTarget] = useState<Topic | null>(null);
  
  useEffect(() => {
    if(!open) return;

    setLoading(true)
    getAllTopics()
      .then(setTopics)
      .catch((err) => {
        console.error("Failed to load topics: ", err);
      })
      .finally(() => setLoading(false))
  }, [open])

  async function handleSaveEdit(topicId: string) {
    const trimmed = editValue.trim()
    if(trimmed === "") { setEditingId(null); return;}
    const original = topics.find(t => t.id === topicId)
    if(original && trimmed === original.label) { setEditingId(null); return;}

    try {
      const updatedTopic = await updateTopic(topicId, editValue)
      
      setTopics(prev => prev.map(t => t.id === topicId ? updatedTopic : t))
    } catch (err) {
      console.error("Failed to update topic: ", err);
    } finally {
      setEditingId(null);
    }
  }

  async function handleConfirmDelete() {
    if (!deleteTarget) return;
    try {
      await deleteTopic(deleteTarget.id); // Server Action, cascades server-side
      setTopics(prev => prev.filter(t => t.id !== deleteTarget.id));
    } catch (err) {
      console.error("Failed to delete topic:", err);
    } finally {
      setDeleteTarget(null);
    }
  }

  return (
    <>
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent
          onEscapeKeyDown={(e) => {
            if (editingId !== null) {
              e.preventDefault();
              setEditingId(null);
            }
          }}
        >
          <DialogHeader>
            <DialogTitle>Edit Topics</DialogTitle>
          </DialogHeader>
          {loading ? <p>Loading...</p> : topics.map(t => (
            <div key={t.id}
            className="group flex items-center justify-between rounded px-2 py-1 hover:bg-muted"
            >
              {editingId === t.id ? (
                <input 
                  autoFocus
                  value={editValue}
                  onChange={(e) => setEditValue(e.target.value) }
                  onKeyDown={(e) => {
                    if (e.key === "Enter") handleSaveEdit(t.id)
                    if (e.key ==="Escape") {
                      e.stopPropagation(); 
                      setEditingId(null);
                    }
                  }}
                  onBlur={() => handleSaveEdit(t.id)}
                  className="bg-transparent outline-none border-b border-primary"
                />
              ) : (
                <span
                onClick={() => {
                  setEditingId(t.id)
                  setEditValue(t.label)
                }}
                className="cursor-text"
                >
                  {t.label}
                </span>
                
              )}
              <button
                onClick={() => setDeleteTarget(t)}
                className="opacity-0 group-hover:opacity-100 transition-opacity"
              >
                <Trash className="h-4 w-4" />
              </button>
            </div>
          ))}
        </DialogContent>
      </Dialog>
      <Dialog
        open={deleteTarget !== null} onOpenChange={(open) => !open && setDeleteTarget(null)}
        >
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Delete &quot;{deleteTarget?.label}&quot;?</DialogTitle>
              <DialogDescription>
                This will remove this topic from all associated grants and clients. This can&apos;t be undone.
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <Button variant="outline" onClick={() => setDeleteTarget(null)}>Cancel</Button>
              <Button variant="destructive" onClick={handleConfirmDelete}>Delete</Button>
            </DialogFooter>
          </DialogContent>
      </Dialog>
    </>
  )
}