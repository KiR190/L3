import { useState } from 'react';
import { useLinkTelegram } from '@/entities/user/queries';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';
import { Alert, AlertDescription } from '@/shared/ui/alert';
import type { ApiError } from '@/shared/types/common';

export const LinkTelegramForm = () => {
  const linkMutation = useLinkTelegram();
  const [username, setUsername] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    
    const cleanUsername = username.startsWith('@') ? username.slice(1) : username;
    
    try {
      await linkMutation.mutateAsync({ telegram_username: cleanUsername });
      setSuccess('Telegram linked! Please send /register to the bot to complete setup.');
      setUsername('');
    } catch (err: unknown) {
      setError((err as ApiError).error || 'Failed to link Telegram');
    }
  };

  return (
    <div className="space-y-4">
      <Alert>
        <AlertDescription>
          <p className="font-semibold mb-2">How to link Telegram:</p>
          <ol className="list-decimal list-inside space-y-1 text-sm">
            <li>Enter your Telegram username below</li>
            <li>Open Telegram and search for the EventBooker bot</li>
            <li>Send <code className="bg-muted px-1 rounded">/register</code> to the bot</li>
            <li>You'll start receiving booking notifications on Telegram!</li>
          </ol>
        </AlertDescription>
      </Alert>
      
      <form onSubmit={handleSubmit} className="space-y-4">
        {error && (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}
        
        {success && (
          <Alert>
            <AlertDescription>{success}</AlertDescription>
          </Alert>
        )}
        
        <div className="space-y-2">
          <Label htmlFor="telegram">Telegram Username</Label>
          <Input
            id="telegram"
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="@username"
            required
          />
        </div>

        <Button type="submit" className="w-full" disabled={linkMutation.isPending}>
          {linkMutation.isPending ? 'Linking...' : 'Link Telegram'}
        </Button>
      </form>
    </div>
  );
};

