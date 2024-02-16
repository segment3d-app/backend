CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE "users" (
    "uid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" VARCHAR(255),
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "avatar" VARCHAR(255),
    "password" VARCHAR(255),
    "provider" VARCHAR(255) NOT NULL,
    "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "passwordChangedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE TABLE "tasks" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "uid" UUID REFERENCES "users"("uid"),
    "title" VARCHAR(255) NOT NULL,
    "status" VARCHAR(255) NOT NULL,
    "url" VARCHAR(255) NOT NULL,
    "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE TABLE "assets" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "uid" UUID REFERENCES "users"("uid"),
    "title" VARCHAR(255) NOT NULL,
    "slug" VARCHAR(255) UNIQUE NOT NULL,
    "assetUrl" VARCHAR(255) NOT NULL,
    "assetType" VARCHAR(255) NOT NULL,
    "thumbnailUrl" VARCHAR(255) NOT NULL,
    "gaussianUrl" VARCHAR(255),
    "pointCloudUrl" VARCHAR(255),
    "isPrivate" BOOLEAN DEFAULT FALSE NOT NULL,
    "status" VARCHAR(255) NOT NULL, -- created, generating splat, completed, failed
    "likes" INT NOT NULL,
    "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE TABLE "tags" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" VARCHAR(255) UNIQUE NOT NULL,
    "slug" VARCHAR(255) UNIQUE NOT NULL,
    "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE TABLE "assetsToTags" (
    "assetsId" UUID REFERENCES "assets"("id"),
    "tagsId" UUID REFERENCES "tags"("id"),
    "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("assetsId", "tagsId")
);
CREATE TABLE "likes" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "uid" UUID REFERENCES "users"("uid"),
    "assetsId" UUID REFERENCES "assets"("id"),
    "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);