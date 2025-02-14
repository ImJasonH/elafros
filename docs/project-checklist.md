# Knative Serving Project Checklist

The purpose of this checklist is to outline must-have items required to open the Knative Serving GitHub repositories to the public.

## Blockers to making project public 

The tasks we would like to complete prior to making repository public, are:

- [x] Migrate code to github.com/knative organization  
- [x] Migrate project roles document to the CONTRIBUTING.md
- [x] Review project roles & governance changes
- [x] Knative Serving design/API documents posted to repo and reviewed 
- [x] Implement Service CRD in the compute/serving repo
- [ ] Seed initial project roadmap docs (e.g. API, Build, Autoscaling, Eventing)
- [ ] Validate no-build "Knative Serving Easy" deployment (aka [install](https://github.com/knative/install)) for GKE and non-GKE k8s cluster
- [ ] Validate all sample apps against knative/install deployed cluster
- [ ] Operator/Developer upgrade user guide (fully manual ok for now)



## After project is public 

Additional items, that while still important can be completed after opening the repository, are:

- [ ] Configure @knative Twitter account profile 
- [ ] Set up long-term billing for Slack, Github, etc.
- [ ] Configure social Stack Overflow label(s)
- [ ] Clean up GitHub landing page (pin repos, update description & links, seed community)
- [ ] Team Drive and mailing lists created
- [ ] Finalize logo design and apply it in Slack/GitHub
- [ ] Migrate to public prow
- [ ] Potentially rename knative/serving to knative/serving or knative/compute
