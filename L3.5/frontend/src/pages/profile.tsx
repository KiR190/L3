import { useState } from 'react';
import { useUser, useUpdateProfile } from '@/entities/user/queries';
import { LinkTelegramForm } from '@/features/telegram/link-telegram-form/link-telegram-form';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/shared/ui/card';
import { Badge } from '@/shared/ui/badge';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';
import { Alert, AlertDescription } from '@/shared/ui/alert';
import type { ApiError } from '@/shared/types/common';

export function ProfilePage() {
  const { data: user, isLoading, error } = useUser();
  const updateProfileMutation = useUpdateProfile();
  
  const [isEditingEmail, setIsEditingEmail] = useState(false);
  const [email, setEmail] = useState('');
  const [updateError, setUpdateError] = useState('');
  const [updateSuccess, setUpdateSuccess] = useState('');

  const handleEmailEdit = () => {
    setEmail(user?.email || '');
    setIsEditingEmail(true);
    setUpdateError('');
    setUpdateSuccess('');
  };

  const handleEmailSave = async () => {
    if (!email || email === user?.email) {
      setIsEditingEmail(false);
      return;
    }

    setUpdateError('');
    setUpdateSuccess('');

    try {
      await updateProfileMutation.mutateAsync({ email });
      setUpdateSuccess('Email updated successfully!');
      setIsEditingEmail(false);
    } catch (err: unknown) {
      setUpdateError((err as ApiError).error || 'Failed to update email');
    }
  };

  const handleNotificationChange = async (preference: 'email' | 'telegram') => {
    if (preference === user?.preferred_notification) return;

    setUpdateError('');
    setUpdateSuccess('');

    try {
      await updateProfileMutation.mutateAsync({ preferred_notification: preference });
      setUpdateSuccess('Notification preference updated!');
    } catch (err: unknown) {
      setUpdateError((err as ApiError).error || 'Failed to update preference');
    }
  };

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-8 max-w-2xl">
        <div className="h-8 bg-muted animate-pulse rounded mb-4"></div>
        <div className="h-64 bg-muted animate-pulse rounded"></div>
      </div>
    );
  }

  if (error || !user) {
    return (
      <Alert variant="destructive">
        <AlertDescription>Failed to load profile</AlertDescription>
      </Alert>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      <h1 className="text-4xl font-bold mb-8">Profile</h1>

      {updateError && (
        <Alert variant="destructive" className="mb-4">
          <AlertDescription>{updateError}</AlertDescription>
        </Alert>
      )}

      {updateSuccess && (
        <Alert className="mb-4">
          <AlertDescription>{updateSuccess}</AlertDescription>
        </Alert>
      )}

      <div className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>Account Information</CardTitle>
            <CardDescription>Your personal details</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <Label className="text-sm text-muted-foreground">User ID</Label>
              <div className="font-mono text-xs mt-1">{user.id}</div>
            </div>

            <div>
              <Label className="text-sm text-muted-foreground">Role</Label>
              <div className="mt-1">
                <Badge variant={user.role === 'admin' ? 'destructive' : 'secondary'}>
                  {user.role === 'admin' ? '👑 Admin' : '👤 User'}
                </Badge>
              </div>
            </div>

            <div>
              <Label className="text-sm text-muted-foreground">Email</Label>
              {isEditingEmail ? (
                <div className="flex gap-2 mt-1">
                  <Input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="your@email.com"
                  />
                  <Button onClick={handleEmailSave} disabled={updateProfileMutation.isPending}>
                    Save
                  </Button>
                  <Button variant="outline" onClick={() => setIsEditingEmail(false)}>
                    Cancel
                  </Button>
                </div>
              ) : (
                <div className="flex items-center gap-2 mt-1">
                  <div className="font-semibold">{user.email}</div>
                  <Button variant="ghost" size="sm" onClick={handleEmailEdit}>
                    Edit
                  </Button>
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Notification Preferences</CardTitle>
            <CardDescription>Choose how you want to receive booking updates</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <Label className="text-sm text-muted-foreground mb-2 block">Preferred Channel</Label>
              <div className="flex gap-2">
                <Button
                  variant={user.preferred_notification === 'email' ? 'default' : 'outline'}
                  onClick={() => handleNotificationChange('email')}
                  disabled={updateProfileMutation.isPending}
                >
                  📧 Email
                </Button>
                <Button
                  variant={user.preferred_notification === 'telegram' ? 'default' : 'outline'}
                  onClick={() => handleNotificationChange('telegram')}
                  disabled={updateProfileMutation.isPending || !user.telegram_registered}
                >
                  💬 Telegram
                </Button>
              </div>
              {user.preferred_notification === 'telegram' && !user.telegram_registered && (
                <p className="text-xs text-muted-foreground mt-2">
                  Link your Telegram account below to enable Telegram notifications
                </p>
              )}
            </div>

            <div className="pt-2">
              <div className="text-sm font-medium mb-1">Current Preference</div>
              <Badge variant={user.preferred_notification === 'telegram' ? 'default' : 'secondary'}>
                {user.preferred_notification === 'telegram' ? '💬 Telegram' : '📧 Email'}
              </Badge>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Telegram Integration</CardTitle>
            <CardDescription>
              {user.telegram_username 
                ? `Connected as @${user.telegram_username}`
                : 'Link your Telegram account for notifications'
              }
            </CardDescription>
          </CardHeader>
          <CardContent>
            {user.telegram_registered && (
              <Alert className="mb-4">
                <AlertDescription>
                  ✅ Your Telegram account is connected and verified!
                </AlertDescription>
              </Alert>
            )}
            
            {!user.telegram_registered && (
              <LinkTelegramForm />
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
