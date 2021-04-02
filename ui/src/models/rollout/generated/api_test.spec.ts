/**
 * pkg/apiclient/rollout/rollout.proto
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * OpenAPI spec version: version not set
 * 
 *
 * NOTE: This file is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 * Do not edit the file manually.
 */

import * as api from "./api"
import { Configuration } from "./configuration"

const config: Configuration = {}

describe("RolloutServiceApi", () => {
  let instance: api.RolloutServiceApi
  beforeEach(function() {
    instance = new api.RolloutServiceApi(config)
  });

  test("rolloutServiceAbortRollout", () => {
    const name: string = "name_example"
    return expect(instance.rolloutServiceAbortRollout(name, {})).resolves.toBe(null)
  })
  test("rolloutServiceGetNamespace", () => {
    return expect(instance.rolloutServiceGetNamespace({})).resolves.toBe(null)
  })
  test("rolloutServiceGetRollout", () => {
    const name: string = "name_example"
    return expect(instance.rolloutServiceGetRollout(name, {})).resolves.toBe(null)
  })
  test("rolloutServiceListRollouts", () => {
    return expect(instance.rolloutServiceListRollouts({})).resolves.toBe(null)
  })
  test("rolloutServicePromoteRollout", () => {
    const name: string = "name_example"
    return expect(instance.rolloutServicePromoteRollout(name, {})).resolves.toBe(null)
  })
  test("rolloutServiceRestartRollout", () => {
    const name: string = "name_example"
    return expect(instance.rolloutServiceRestartRollout(name, {})).resolves.toBe(null)
  })
  test("rolloutServiceRetryRollout", () => {
    const name: string = "name_example"
    return expect(instance.rolloutServiceRetryRollout(name, {})).resolves.toBe(null)
  })
  test("rolloutServiceSetRolloutImage", () => {
    const rollout: string = "rollout_example"
    const container: string = "container_example"
    const image: string = "image_example"
    const tag: string = "tag_example"
    return expect(instance.rolloutServiceSetRolloutImage(rollout, container, image, tag, {})).resolves.toBe(null)
  })
  test("rolloutServiceUndoRollout", () => {
    const rollout: string = "rollout_example"
    const revision: string = "revision_example"
    return expect(instance.rolloutServiceUndoRollout(rollout, revision, {})).resolves.toBe(null)
  })
  test("rolloutServiceVersion", () => {
    return expect(instance.rolloutServiceVersion({})).resolves.toBe(null)
  })
  test("rolloutServiceWatchRollout", () => {
    const name: string = "name_example"
    return expect(instance.rolloutServiceWatchRollout(name, {})).resolves.toBe(null)
  })
  test("rolloutServiceWatchRollouts", () => {
    return expect(instance.rolloutServiceWatchRollouts({})).resolves.toBe(null)
  })
})

