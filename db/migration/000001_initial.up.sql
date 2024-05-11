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
CREATE TABLE "assets" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "uid" UUID REFERENCES "users"("uid") NOT NULL,
    "title" VARCHAR(255) NOT NULL,
    "slug" VARCHAR(255) UNIQUE NOT NULL,
    "type" VARCHAR(255) NOT NULL,
    "thumbnailUrl" VARCHAR(255) NOT NULL,
    "photoDirUrl" VARCHAR(255) NOT NULL,
    "splatUrl" VARCHAR(255),
    "pclUrl" VARCHAR(255),
    "pclColmapUrl" VARCHAR(255),
    "segmentedPclDirUrl" VARCHAR(255),
    "segmentedSplatDirUrl" VARCHAR(255),
    "isPrivate" BOOLEAN DEFAULT FALSE NOT NULL,
    "status" VARCHAR(255) NOT NULL, -- created, generating colmap,  generating splat, completed, failed
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
    "assetsId" UUID NOT NULL,
    "tagsId" UUID REFERENCES "tags"("id") NOT NULL,
    "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY ("assetsId") REFERENCES "assets"("id") ON DELETE CASCADE,
    PRIMARY KEY ("assetsId", "tagsId")
);
CREATE TABLE "likes" (
    "uid" UUID NOT NULL REFERENCES "users"("uid"),
    "assetsId" UUID NOT NULL,
    "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY ("assetsId") REFERENCES "assets"("id") ON DELETE CASCADE,
    PRIMARY KEY ("uid", "assetsId")
);
