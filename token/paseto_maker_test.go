package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xiaowuzai/simplebank/util"
)

func TestPasetoMaker(t *testing.T) {
	// 测试JWTMaker的New方法
	symmetricKey := util.NewPasetoSymmetricKey()
	// symmetricKey := "a0bb8951f8aa97d94707503d7560478ee1bd2db82868bee11113ca1317aac098"

	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomOwner()
	duration := time.Minute
	// 测试JWTMaker的Create方法
	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, token)

	// 测试JWTMaker的Verify方法
	payloadV, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payloadV)
	require.NotEmpty(t, payloadV.ID)

	require.Equal(t, payload.Username, payloadV.Username)
	require.WithinDuration(t, payload.ExpiredAt, payloadV.ExpiredAt, time.Second)
	require.WithinDuration(t, payload.IssuedAt, payloadV.IssuedAt, time.Second)
}

func TestVerifyPasetoToken(t *testing.T) {
	symmetricKey := util.NewPasetoSymmetricKey()

	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	token, payload, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, token)

	payloadV, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	// require.True(t, strings.Contains(err.Error(), ErrExpiredToken.Error())) // TODO: 如何断言错误？
	require.Nil(t, payloadV)
}
