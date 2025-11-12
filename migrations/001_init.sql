-- Migration 001: Initial schema for newsletter scheduler system

-- Topics table: stores newsletter topics that admins can create
CREATE TABLE topics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Subscribers table: stores email subscribers
CREATE TABLE subscribers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Subscriptions table: many-to-many relationship between subscribers and topics
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscriber_id UUID NOT NULL REFERENCES subscribers(id) ON DELETE CASCADE,
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    subscribed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_active BOOLEAN DEFAULT true,
    UNIQUE(subscriber_id, topic_id)
);

-- Content table: stores scheduled newsletter content
CREATE TABLE content (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    subject VARCHAR(500) NOT NULL,
    body TEXT NOT NULL,
    send_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(50) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'sent', 'failed', 'cancelled')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Deliveries table: tracks individual email deliveries
CREATE TABLE deliveries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    subscriber_id UUID NOT NULL REFERENCES subscribers(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'failed', 'bounced')),
    sent_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(content_id, subscriber_id)
);

-- Job scheduler table: durable job scheduling
CREATE TABLE job_scheduler (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    job_type VARCHAR(100) NOT NULL,
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'enqueued', 'completed', 'failed')),
    attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for better performance
CREATE INDEX idx_topics_name ON topics(name);
CREATE INDEX idx_subscribers_email ON subscribers(email);
CREATE INDEX idx_subscriptions_subscriber_id ON subscriptions(subscriber_id);
CREATE INDEX idx_subscriptions_topic_id ON subscriptions(topic_id);
CREATE INDEX idx_content_topic_id ON content(topic_id);
CREATE INDEX idx_content_send_at ON content(send_at);
CREATE INDEX idx_content_status ON content(status);
CREATE INDEX idx_deliveries_content_id ON deliveries(content_id);
CREATE INDEX idx_deliveries_subscriber_id ON deliveries(subscriber_id);
CREATE INDEX idx_deliveries_status ON deliveries(status);
CREATE INDEX idx_job_scheduler_scheduled_at ON job_scheduler(scheduled_at);
CREATE INDEX idx_job_scheduler_status ON job_scheduler(status);

-- Updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply updated_at triggers
CREATE TRIGGER update_topics_updated_at BEFORE UPDATE ON topics FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_subscribers_updated_at BEFORE UPDATE ON subscribers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_content_updated_at BEFORE UPDATE ON content FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_deliveries_updated_at BEFORE UPDATE ON deliveries FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_job_scheduler_updated_at BEFORE UPDATE ON job_scheduler FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
